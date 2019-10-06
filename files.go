package main

import (
	"errors"
	"io"
	"strings"

	"github.com/fatih/color"
)

type Type string

const (
	TypeFile   Type = "FILE"
	TypeFolder Type = "FOLDER"
)

var (
	folderColor = color.New(color.FgBlue)
	symColor    = color.New(color.FgHiCyan)
)

type FileInfo interface {
	Name() string
	Type() Type
	Path() string
	IsLast() bool
	Write(w io.Writer, isFullPath bool) error
}

type File struct {
	name   string
	parent FileInfo
	isLast bool

	path    string
	SymLink string
}

func NewFile(name string, parent FileInfo, isLast bool, symLink string) *File {
	return &File{
		name:    name,
		parent:  parent,
		isLast:  isLast,
		SymLink: symLink,
	}
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Suffix() string {
	n := strings.Split(f.name, ".")
	return n[len(n)-1]
}

func (f *File) Type() Type {
	return TypeFile
}

func (f *File) Path() string {
	if f.path == "" {
		if f.parent == nil {
			f.path = f.name
		} else {
			f.path = f.parent.Path() + "/" + f.name
		}
	}

	return f.path
}

func (f *File) IsLast() bool {
	return f.isLast
}

func (f *File) IsSym() bool {
	return f.SymLink != ""
}

func (f *File) Write(w io.Writer, isFullPath bool) error {
	icon, ok := icons[f.Suffix()]
	if !ok {
		icon = defaultFileIcon
	}

	switch v := f.parent.(type) {
	case *Folder:
		if f.IsLast() {
			w.Write([]byte(v.ChildPrefix + "└── "))
		} else {
			w.Write([]byte(v.ChildPrefix + "├── "))
		}

		if f.IsSym() {
			color.New(icon.Color).Fprint(w, icon.Icon+" ")
			if isFullPath {
				symColor.Fprint(w, f.Path())
			} else {
				symColor.Fprint(w, f.Name())
			}

			w.Write([]byte(" -> " + f.SymLink + "\n"))
		} else {
			color.New(icon.Color).Fprint(w, icon.Icon+" ")
			if isFullPath {
				w.Write([]byte(f.Path() + "\n"))
			} else {
				w.Write([]byte(f.Name() + "\n"))
			}
		}
	default:
		return errors.New("Unexpected parent type")
	}
	return nil
}

type Folder struct {
	name        string
	parent      FileInfo
	Children    []FileInfo
	ChildPrefix string
	isLast      bool

	path string
}

func NewFolder(name string, parent FileInfo, isLast bool) *Folder {
	return &Folder{
		name:   name,
		parent: parent,
		isLast: isLast,
	}
}

func (f *Folder) Name() string {
	return f.name
}

func (f *Folder) Type() Type {
	return TypeFolder
}

func (f *Folder) Path() string {
	if f.path == "" {
		if f.parent == nil {
			f.path = f.name
		} else {
			f.path = f.parent.Path() + "/" + f.name
		}
	}

	return f.path
}

func (f *Folder) IsLast() bool {
	return f.isLast
}

func (f *Folder) Write(w io.Writer, isFullPath bool) error {
	if f.parent == nil {
		folderColor.Fprintln(w, f.Name())
		return nil
	}

	switch v := f.parent.(type) {
	case *Folder:
		if f.IsLast() {
			w.Write([]byte(v.ChildPrefix + "└── "))
			f.ChildPrefix = v.ChildPrefix + "    "
		} else {
			w.Write([]byte(v.ChildPrefix + "├── "))
			f.ChildPrefix = v.ChildPrefix + "│   "
		}

		if isFullPath {
			folderColor.Fprintln(w, f.Path())
		} else {
			folderColor.Fprintln(w, f.Name())
		}
	default:
		return errors.New("Unexpected parent type")
	}
	return nil
}
