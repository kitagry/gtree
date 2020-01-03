package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

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
	All []bool `short:"a" long:"all" description:"All files are listed."`

	OnlyDirectory []bool `short:"d" description:"List directories only."`

	IgnorePattern string `short:"I" description:"Do not list files that match the given pattern."`
}

// IsAll returns true, if user specify '-a' or '-all' option.
func (l *ListSearchOptions) IsAll() bool {
	return len(l.All) != 0
}

// IsOnlyDirectry returns true, if user specify '-d' option.
func (l *ListSearchOptions) IsOnlyDirectry() bool {
	return len(l.OnlyDirectory) != 0
}

// ListDisplayOptions is options which use when display file tree.
type ListDisplayOptions struct {
	FullPath []bool `short:"f" description:"Print the full path prefix for each file."`

	Output string `short:"o" description:"Output to file instead of stdout."`
}

// IsFullPath returns true, if user specify '-f' option.
func (l *ListDisplayOptions) IsFullPath() bool {
	return len(l.FullPath) != 0
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
		fmt.Println("gtree v0.2")
		os.Exit(0)
	}

	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gtree"
	parser.Usage = "[-adf] [--version] [-I pattern] [-o filename] [--help] [--] [<directory list>]"
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
	var out io.Writer
	if outputFile := opts.ListOptions.ListDisplayOptions.Output; outputFile != "" {
		var err error
		if _, err = os.Stat(outputFile); !os.IsNotExist(err) {
			fmt.Printf("Output file already exists. Are you sure to overwrite %s?[Y/n] ", outputFile)

			var answer string
			fmt.Scan(&answer)
			if strings.ToLower(strings.TrimRight(answer, "\n")) != "y" {
				return fmt.Errorf("output file already exists")
			}
		}

		out, err = os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("file create/open error: %v", err)
		}
		defer out.(*os.File).Close()
	} else {
		out = os.Stdout
	}

	w := bufio.NewWriter(out)
	for file := range ch {
		err := file.Write(w, opts.ListOptions.ListDisplayOptions.IsFullPath())
		if err != nil {
			w.Flush()
			return err
		}
	}
	w.Flush()
	return nil
}
