package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	IsAll []bool `short:"a" description:"All files are listed."`
}

var (
	opts Options

	folderColor = color.New(color.FgCyan)
	symColor    = color.New(color.FgHiCyan)
)

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gtree"
	parser.Usage = "gtree [-a] [<directory list>]"

	directories, err := parser.Parse()
	if err != nil {
		panic(err)
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
	err := Dirwalk(&rootFile)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(os.Stdout)
	_ = rootFile.Write(w)
	w.Flush()
	return nil
}

func Dirwalk(root FileInfo) error {
	switch f := root.(type) {
	case *File:
		return nil
	case *Folder:
		files, err := ioutil.ReadDir(f.Path())
		if err != nil {
			return err
		}

		for i, file := range files {
			if len(opts.IsAll) == 0 && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			isLast := i == len(files)-1
			if file.IsDir() {
				child := NewFolder(file.Name(), f, isLast)
				Dirwalk(&child)
				f.Children = append(f.Children, &child)
				continue
			}

			child := NewFile(file.Name(), f, isLast)
			if file.Mode()&os.ModeSymlink != 0 {
				sym, _ := os.Readlink(f.Path() + "/" + file.Name())
				child.SymLink = sym
			}
			f.Children = append(f.Children, &child)
		}
	default:
		return errors.New("Unexpected File type")
	}
	return nil
}
