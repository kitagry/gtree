package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"golang.org/x/xerrors"
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

	IgnorePatterns []string `short:"I" description:"Do not list files that match the given pattern."`

	Level *int `short:"L" long:"level" description:"Descend only level directories deep."`
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

	NoIcons []bool `short:"n" description:"Do not show the icon of files and directories"`
}

// IsFullPath returns true, if user specify '-f' option.
func (l *ListDisplayOptions) IsFullPath() bool {
	return len(l.FullPath) != 0
}

// NoIcon returns true, if user specify '-n' option.
func (l *ListDisplayOptions) NoIcon() bool {
	return len(l.NoIcons) != 0
}

type MiscellaneousOptions struct {
	Version func() `long:"version" description:"show version"`
}

// Options is all options.
type Options struct {
	ListOptions          *ListOptions          `group:"List Options"`
	MiscellaneousOptions *MiscellaneousOptions `group:"Miscellaneous Options"`
}

func newOptionsParser(opts *Options) *flags.Parser {
	opts.ListOptions = &ListOptions{}
	opts.MiscellaneousOptions = &MiscellaneousOptions{}

	opts.MiscellaneousOptions.Version = func() {
		fmt.Println("gtree v0.2")
		os.Exit(0)
	}

	parser := flags.NewParser(opts, flags.Default)
	parser.Name = "gtree"
	parser.Usage = "[-adfn] [--version] [-I pattern] [-o filename] [-L level] [--help] [--] [<directory list>]"
	return parser
}

// run is to search files, and display these files.
// Search and Display communicate through channel.
func run(ctx context.Context) int {
	var opts Options
	parser := newOptionsParser(&opts)

	directories, err := parser.Parse()
	if err != nil {
		return statusErr
	}

	if len(directories) == 0 {
		directories = append(directories, ".")
	}

	for _, d := range directories {
		err = showTree(d, opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", parser.Name, err)
			return statusErr
		}
	}

	return statusOK
}

func showTree(root string, opts Options) error {
	if opts.ListOptions.ListSearchOptions.Level != nil && *opts.ListOptions.ListSearchOptions.Level <= 0 {
		return fmt.Errorf("Invalid level, must be greater than 0.")
	}

	f, err := os.Stat(root)
	if err != nil {
		return xerrors.Errorf("failed to find root: %v", err)
	}

	base, _ := filepath.Split(root)
	rootFile := NewFileInfoForBase(f, nil, base, true)
	ch := make(chan FileInfo)

	if !rootFile.IsDir() {
		errRootIsNotDir := fmt.Errorf("%s is not dir", rootFile.Name())
		rootFile.SetError(errRootIsNotDir)
	}

	// Search files.
	go Dirwalk(rootFile, ch, opts.ListOptions.ListSearchOptions)

	// Display files.
	var out io.Writer

	if outputFile := opts.ListOptions.ListDisplayOptions.Output; outputFile != "" {
		var err error
		if err = checkOverWrite(outputFile); err != nil {
			return xerrors.Errorf("denided overwrite: %w", err)
		}

		out, err = os.Create(outputFile)
		if err != nil {
			return xerrors.Errorf("file create/open error: %w", err)
		}
		defer out.(*os.File).Close()
	} else {
		out = os.Stdout
	}

	w := bufio.NewWriter(out)
	p := NewPrinter(opts.ListOptions.ListDisplayOptions)

	for file := range ch {
		err := p.Write(w, file)
		if err != nil {
			return err
		}
	}

	w.Flush()
	return nil
}

var errFileExist = fmt.Errorf("output file already exists")

func checkOverWrite(filename string) error {
	var err error
	if _, err = os.Stat(filename); !os.IsNotExist(err) {
		fmt.Printf("Output file already exists. Are you sure to overwrite %s?[Y/n] ", filename)

		var answer string
		if _, err = fmt.Scan(&answer); err != nil {
			return xerrors.Errorf("failed to scan answer: %v", err)
		}

		if strings.ToLower(strings.TrimRight(answer, "\n")) != "y" {
			return errFileExist
		}
	}
	return nil
}
