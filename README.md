# gtree

🎄 tree command with icons.

![sample image](./images/sample.png)

## Requirements

[Nerd Fonts](https://www.nerdfonts.com/) or related fonts.

## Installation

```
$ go get github.com/kitagry/gtree
```

## Usage

```
$ gtree -h
Usage:
  gtree [-adfn] [--version] [-I pattern] [-o filename] [-L level] [--help] [--]
[<directory list>]

List Options:
  -a, --all      All files are listed.
  -d             List directories only.
  -I=            Do not list files that match the given pattern.
  -L, --level=   Descend only level directories deep.
  -f             Print the full path prefix for each file.
  -o=            Output to file instead of stdout.
  -n             Do not show the icon of files and directories

Miscellaneous Options:
      --version  show version

Help Options:
  -h, --help     Show this help message
```
