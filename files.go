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

type FileInfo interface {
	Name() string
	Type() Type
	Path() string
	IsLast() bool
	Write(w io.Writer) error
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

func (f *File) Write(w io.Writer) error {
	icon, ok := icons[f.Suffix()]
	if !ok {
		icon = defaultFileIcon
	}

	switch v := f.parent.(type) {
	case *Folder:
		if f.IsSym() {
			if f.IsLast() {
				w.Write([]byte(v.ChildPrefix + "└── "))
			} else {
				w.Write([]byte(v.ChildPrefix + "├── "))
			}
			color.New(icon.Color).Fprint(w, icon.Icon+" ")
			symColor.Fprint(w, f.Name())
			w.Write([]byte(" -> " + f.SymLink + "\n"))
		} else {
			if f.IsLast() {
				w.Write([]byte(v.ChildPrefix + "└── "))
				color.New(icon.Color).Fprint(w, icon.Icon+" ")
				w.Write([]byte(f.Name() + "\n"))
			} else {
				w.Write([]byte(v.ChildPrefix + "├── "))
				color.New(icon.Color).Fprint(w, icon.Icon+" ")
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

func (f *Folder) Write(w io.Writer) error {
	if f.parent == nil {
		w.Write([]byte(f.Name() + "\n"))
		return nil
	}

	switch v := f.parent.(type) {
	case *Folder:
		if f.IsLast() {
			w.Write([]byte(v.ChildPrefix + "└── "))
			folderColor.Fprintln(w, f.Name())
			f.ChildPrefix = v.ChildPrefix + "    "
		} else {
			w.Write([]byte(v.ChildPrefix + "├── "))
			folderColor.Fprintln(w, f.Name())
			f.ChildPrefix = v.ChildPrefix + "│   "
		}
	default:
		return errors.New("Unexpected parent type")
	}
	return nil
}
