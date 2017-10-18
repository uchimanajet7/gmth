# gmth
A command line tool converts Markdown file to HTML file.

## Description
You can convert Markdown file to HTML file.

You can execute arbitrary commands before starting the conversion process and after the processing has finished.

You can execute character string replacement processing on the converted HTML file.


## Features
- It is made by golang so it supports multi form.
- You can control the operation in the setting file.
	- You can create a file that can be displayed as an HTML page by specifying CSS.
	- You can execute arbitrary commands before starting conversion process and after processing.
	- You can execute character string replacement processing on the converted HTML file.


## Requirement
- Go 1.9+
- Packages in use
	- russross/blackfriday: Blackfriday: a markdown processor for Go
		- https://github.com/russross/blackfriday
	- fsnotify/fsnotify: Cross-platform file system notifications for Go.
		- https://github.com/fsnotify/fsnotify
	- spf13/cobra: A Commander for modern Go CLI interactions
		- https://github.com/spf13/cobra

## Usage
Just by specifying the markdown file path in the run command.

```	sh
$ ./gmth run -f markdown.md
```

However, setting is necessary to execute.

### Setting Example

1. In the same place as the binary file create execution settings file.

1. Execution settings are done with `config.json` file.

```json
{
    "Page": true,
    "CSS": "style.css",
    "PreCommands": [["textlint", "%INPUT_PATH%"]],
	"PostCommands": [["open", "%OUTPUT_PAGE_PATH%"]],
	"ReplaceTexts": [
        "<blockquote",
        "<blockquote class=\"is-colored\"",
        "<ol",
        "<ol class=\"list-colored\"",
        "<img",
        "<img class=\"image-with-border\"",
        "<table",
        "<table class=\"table table--bordered table--color-header\""
    ]
}
```

- About setting items
	- `Page`: Boolean
		- Whether to output a complete HTML file.
	- `CSS`: String
		- Valid when `Page` item is **true**.
		- Specify the path of the CSS file to be used for the HTML file.
	- `PreCommands`: Array
		- Multidimensional array.
		- You can specify the command you want to execute before conversion starts.
		- One command with one array.
		- Arrays are pipelined and executed.
		- `%INPUT_PATH%` is replaced with the path of the input file.
	- `PostCommands `: Array
		- Multidimensional array.
		- You can specify the command to be executed after conversion is completed.
		- One command with one array.
		- Arrays are pipelined and executed.
		- `%OUTPUT_PATH%` is replaced with the path of the output file.
		- `%OUTPUT_PAGE_PATH%` is replaced with the path of the output file.
			- Valid when `Page` item is **true**.
	- `ReplaceTexts`: Array
		- String List a list to replace.

Current specification is a specification that continues execution even if there is an error in command execution before / after change.

More information, refer to the help command.


```sh
$ ./gmth help
A command line tool converts Markdown file to HTML file

Usage:
  gmth [command]

Available Commands:
  help        Help about any command
  run         Markdown Convert the file to HTML file once
  version     Print the version number of gmth
  watch       Watch Markdown file and converts it to HTML file

Flags:
  -f, --file string   Specify markdown file. [required]
  -h, --help          help for gmth

Use "gmth [command] --help" for more information about a command.
```

## Installation

Please select the package file for your own environment from the releases page, download and unpack it, and put the executable file in a place where included in PATH.

- Releases · uchimanajet7/ca-cli
	- https://github.com/uchimanajet7/gmth/releases

If you build from source yourself.

```	console
$ go get github.com/uchimanajet7/gmth
$ cd $GOPATH/src/github.com/uchimanajet7/gmth
$ make
```

## Author
[uchimanajet7](https://github.com/uchimanajet7)


## Licence
[MIT License](https://github.com/uchimanajet7/gmth/blob/master/LICENSE)
