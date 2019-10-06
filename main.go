package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

type ListOptions struct {
	IsAll []bool `short:"a" long:"all" description:"All files are listed."`

	OnlyDirectory []bool `short:"d" description:"List directories only."`

	FullPath []bool `short:"f" description:"Print the full path prefix for each file."`
}

type MiscellaneousOptions struct {
	Version func() `long:"version" description:"show version"`
}

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
	parser.Usage = "[-ad] [--version] [--] [<directory list>]"
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

func Run(root string) error {
	rootFile := NewFolder(root, nil, true)
	ch := make(chan FileInfo)
	go Dirwalk(rootFile, ch)

	w := bufio.NewWriter(os.Stdout)
	for file := range ch {
		err := file.Write(w, len(opts.ListOptions.FullPath) != 0)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
