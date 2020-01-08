package main

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

var (
	folderColor = color.New(color.FgBlue)
	symColor    = color.New(color.FgHiCyan)
)

// Printer write FileInfo as tree
type Printer struct {
}

// NewPrinter return Printer pointer
func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Write(w io.Writer, f FileInfo, isFullPath bool) error {
	if pa, ok := f.Parent(); ok {
		if f.IsLast() {
			w.Write([]byte(pa.ChildPrefix() + "└── "))
		} else {
			w.Write([]byte(pa.ChildPrefix() + "├── "))
		}
	}

	if !f.IsDir() {
		w.Write([]byte(NewIconString(f.FileType()) + " "))
	}

	var writtenName string
	if isFullPath {
		writtenName = f.Path()
	} else {
		writtenName = f.Name()
	}

	if f.IsDir() {
		w.Write([]byte(folderColor.Sprintf(writtenName)))
	} else if f.IsSym() {
		symLink, err := f.SymLink()
		if err != nil {
			return err
		}

		w.Write([]byte(fmt.Sprintf("%s -> %s", symColor.Sprint(writtenName), symLink)))
	} else {
		w.Write([]byte(writtenName))
	}

	if err := f.Error(); err != nil {
		w.Write([]byte(fmt.Sprintf(" [%s]", err.Error())))
	}

	w.Write([]byte("\n"))
	return nil
}
