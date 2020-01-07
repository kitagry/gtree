package main

import (
	"fmt"
	"io"
)

type Printer struct {
}

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

	if f.IsSym() {
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
