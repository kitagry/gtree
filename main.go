package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

// ListOptions is Options for basic gtree command.
type ListOptions struct {
	// ListSearchOptions is options which use when searching file tree.
	ListSearchOptions *ListSearchOptions

	// ListDisplayOptions is options which use when display file tree.
	ListDisplayOptions *ListDisplayOptions
}

// ListSearchOptions is options which use when searching file tree.
type ListSearchOptions struct {
	IsAll []bool `short:"a" long:"all" description:"All files are listed."`

	OnlyDirectory []bool `short:"d" description:"List directories only."`
}

// ListDisplayOptions is options which use when display file tree.
type ListDisplayOptions struct {
	FullPath []bool `short:"f" description:"Print the full path prefix for each file."`
}

type MiscellaneousOptions struct {
	Version func() `long:"version" description:"show version"`
}

// Options is all options.
type Options struct {
	ListOptions          *ListOptions          `group:"List Options"`
	MiscellaneousOptions *MiscellaneousOptions `group:"Miscellaneous Options"`
}

var (
	opts Options
)

func newOptionsParser(opt *Options) *flags.Parser {
	opt.ListOptions = &ListOptions{}
	opt.MiscellaneousOptions = &MiscellaneousOptions{}

	opts.MiscellaneousOptions.Version = func() {
		fmt.Println("gtree v0.0.0")
		os.Exit(0)
	}

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gtree"
	parser.Usage = "[-adf] [--version] [--help] [--] [<directory list>]"
	return parser
}

func main() {
	parser := newOptionsParser(&opts)

	directories, err := parser.Parse()
	if err != nil {
		return
	}

	if len(directories) == 0 {
		directories = append(directories, ".")
	}

	for _, d := range directories {
		err = Run(d)
		if err != nil {
			panic(err)
		}
	}
}

// Run is to search files, and display these files.
// Search and Display communicate through channel.
func Run(root string) error {
	rootFile := NewFolder(root, nil, true)
	ch := make(chan FileInfo)

	// Search files.
	go Dirwalk(rootFile, ch, opts.ListOptions.ListSearchOptions)

	// Display files.
	w := bufio.NewWriter(os.Stdout)
	for file := range ch {
		err := file.Write(w, len(opts.ListOptions.ListDisplayOptions.FullPath) != 0)
		if err != nil {
			w.Flush()
			return err
		}
	}
	w.Flush()
	return nil
}
