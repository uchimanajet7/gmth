package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/russross/blackfriday"
)

var loadedConfig *userConfig

type userConfig struct {
	Page         bool
	CSS          string
	ReplaceTexts []string
	PreCommands  [][]string
	PostCommands [][]string
}

func convertFile(path string, config *userConfig) error {
	// load markdown file
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// pre commands
	if len(config.PreCommands) > 0 {
		out, err := runCommands(config.PreCommands, strings.NewReplacer("%INPUT_PATH%", path))
		if err != nil {
			fmt.Printf("[PreCommands] %+v\n", err)
		}
		fmt.Printf("[PreCommands] %+v\n", out)
	}

	// output partial HTML
	outPath := getOutputPath(path, false)
	html := createHTML(input, false, config.CSS, config.ReplaceTexts)
	err = saveHTML(outPath, html)
	if err != nil {
		return err
	}

	// output page HTML
	outPagePath := getOutputPath(path, true)
	if config.Page {
		pageHTML := createHTML(input, true, config.CSS, config.ReplaceTexts)
		err = saveHTML(outPagePath, pageHTML)
		if err != nil {
			return err
		}
	}

	// post commands
	if len(config.PostCommands) > 0 {
		out, err := runCommands(config.PostCommands, strings.NewReplacer("%OUTPUT_PATH%", outPath, "%OUTPUT_PAGE_PATH%", outPagePath))
		if err != nil {
			fmt.Printf("[PostCommands] %+v\n", err)
		}
		fmt.Printf("[PostCommands] %+v\n", out)
	}

	return nil
}

func saveHTML(path string, html string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file
	_, err = f.WriteString(html)
	if err != nil {
		return err
	}

	return err
}

func getOutputPath(path string, page bool) string {
	dir, file := filepath.Split(path)
	ext := filepath.Ext(path)

	name := strings.Replace(file, ext, "", -1)

	if page {
		name = name + "_page"
	}
	name = name + ".html"

	return filepath.Join(dir, name)
}

func createHTML(input []byte, page bool, css string, replaces []string) string {
	// set up options same
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_HEADER_IDS
	extensions |= blackfriday.EXTENSION_BACKSLASH_LINE_BREAK
	extensions |= blackfriday.EXTENSION_DEFINITION_LISTS

	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_DASHES
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	if page {
		htmlFlags |= blackfriday.HTML_COMPLETE_PAGE
	}

	renderer := blackfriday.HtmlRenderer(htmlFlags, "Automatically generated page", css)
	output := blackfriday.Markdown(input, renderer, extensions)

	replacer := strings.NewReplacer(replaces...)
	outputText := replacer.Replace(string(output))

	return outputText
}

func runCommands(commands [][]string, replacer *strings.Replacer) (string, error) {
	// replace text
	cmds := make([][]string, len(commands))
	for i, v := range commands {
		cmd := make([]string, len(v))
		for j, m := range v {
			cmd[j] = replacer.Replace(m)
		}
		cmds[i] = cmd
	}
	out, err := runPipeline(cmds...)

	return string(out), err
}

func runPipeline(commands ...[]string) ([]byte, error) {
	// run one command
	if len(commands) <= 1 {
		cmd := commands[0]
		cmdName := cmd[0]
		cmdArgs := cmd[1:]
		return exec.Command(cmdName, cmdArgs...).Output()
	}

	cmdList := make([]*exec.Cmd, len(commands))

	for i, v := range commands {
		cmdList[i] = exec.Command(v[0], v[1:]...)
		if i > 0 {
			r, err := cmdList[i-1].StdoutPipe()
			if err != nil {
				return nil, err
			}
			cmdList[i].Stdin = r
		}
		cmdList[i].Stderr = os.Stderr
	}

	var b bytes.Buffer
	cmdList[len(cmdList)-1].Stdout = &b

	for _, v := range cmdList {
		if err := v.Start(); err != nil {
			return nil, err
		}
	}
	for _, v := range cmdList {
		if err := v.Wait(); err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}

func getExecDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(execPath), nil
}

func getConfigPath() (string, error) {
	dir, err := getExecDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.json"), err
}

func loadConfig() (*userConfig, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &userConfig{}
	err = json.NewDecoder(f).Decode(conf)

	return conf, err
}

func saveConfig(config *userConfig) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// write file
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(config)
	if err != nil {
		return err
	}

	return err
}

func getConfig() (*userConfig, error) {
	if loadedConfig != nil {
		return loadedConfig, nil
	}

	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	loadedConfig = config

	return loadedConfig, nil
}
