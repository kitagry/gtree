package main

import (
	"fmt"
	"io"

	"github.com/gookit/color"
	"golang.org/x/xerrors"
)

var (
	folderColor = color.New(defaultFolderIcon.Color)
	symColor    = color.New(color.FgLightCyan)
)

// Printer write FileInfo as tree
type Printer struct {
}

// NewPrinter return Printer pointer
func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Write(w io.Writer, f FileInfo, isFullPath bool) error {
	var err error
	if pa, ok := f.Parent(); ok {
		if f.IsLast() {
			_, err = w.Write([]byte(pa.ChildPrefix() + "└── "))
		} else {
			_, err = w.Write([]byte(pa.ChildPrefix() + "├── "))
		}

		if err != nil {
			return xerrors.Errorf("failed to write: %w", err)
		}
	}

	if !f.IsDir() {
		_, err = w.Write([]byte(NewIconString(f.FileType()) + " "))

		if err != nil {
			return xerrors.Errorf("failed to write: %w", err)
		}
	}

	var writtenName string
	if isFullPath {
		writtenName = f.Path()
	} else {
		writtenName = f.Name()
	}

	switch {
	case f.IsDir():
		_, err = w.Write([]byte(folderColor.Sprintf("%s %s", defaultFolderIcon.Icon, writtenName)))
	case f.IsSym():
		var symLink string
		symLink, err = f.SymLink()
		if err != nil {
			return xerrors.Errorf("failed to retrieve symlink path: %w", err)
		}

		_, err = w.Write([]byte(fmt.Sprintf("%s -> %s", symColor.Sprint(writtenName), symLink)))
	default:
		_, err = w.Write([]byte(writtenName))
	}

	if err != nil {
		return xerrors.Errorf("failed to write: %w", err)
	}

	if err := f.Error(); err != nil {
		_, err = w.Write([]byte(fmt.Sprintf(" [%s]", err.Error())))
		if err != nil {
			return xerrors.Errorf("failed to write: %w", err)
		}
	}

	_, err = w.Write([]byte("\n"))
	if err != nil {
		return xerrors.Errorf("failed to write: %w", err)
	}
	return nil
}
