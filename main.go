package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

type ListOptions struct {
	IsAll []bool `short:"a" long:"all" description:"All files are listed."`
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
	parser.Usage = "[-a] [--version] [--] [<directory list>]"
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
		err := file.Write(w)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func Dirwalk(root FileInfo, ch chan<- FileInfo) {
	dirwalk(root, ch)
	close(ch)
}

func dirwalk(root FileInfo, ch chan<- FileInfo) {
	ch <- root
	switch f := root.(type) {
	case *File:
		return
	case *Folder:
		files, err := ioutil.ReadDir(f.Path())
		if err != nil {
			fmt.Println(err)
			close(ch)
			return
		}

		for i, file := range files {
			if len(opts.ListOptions.IsAll) == 0 && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			isLast := i == len(files)-1
			var child FileInfo
			if file.IsDir() {
				child = NewFolder(file.Name(), f, isLast)
			} else {
				if file.Mode()&os.ModeSymlink != 0 {
					sym, _ := os.Readlink(f.Path() + "/" + file.Name())
					child = NewFile(file.Name(), f, isLast, sym)
				} else {
					child = NewFile(file.Name(), f, isLast, "")
				}
			}

			f.Children = append(f.Children, child)
			dirwalk(child, ch)
		}
	default:
		fmt.Println("Unexpected File type")
		close(ch)
		return
	}
}
