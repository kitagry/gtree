package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
)

type Type string

const (
	TypeFile   Type = "FILE"
	TypeFolder Type = "FOLDER"
)

var (
	folderColor = color.New(color.FgCyan)
	symColor    = color.New(color.FgHiCyan)
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

func NewFile(name string, parent FileInfo, isLast bool) File {
	return File{
		name:   name,
		parent: parent,
		isLast: isLast,
	}
}

func (f *File) Name() string {
	return f.name
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
	switch v := f.parent.(type) {
	case *Folder:
		if f.IsSym() {
			if f.IsLast() {
				w.Write([]byte(v.ChildPrefix + "└── "))
				symColor.Fprint(w, f.Name())
			} else {
				w.Write([]byte(v.ChildPrefix + "├── "))
			}
			symColor.Fprint(w, f.Name())
			w.Write([]byte(" -> " + f.SymLink + "\n"))
		} else {
			if f.IsLast() {
				w.Write([]byte(v.ChildPrefix + "└── " + f.Name() + "\n"))
			} else {
				w.Write([]byte(v.ChildPrefix + "├── " + f.Name() + "\n"))
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

func NewFolder(name string, parent FileInfo, isLast bool) Folder {
	return Folder{
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
		for _, child := range f.Children {
			err := child.Write(w)
			if err != nil {
				return err
			}
		}
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

	for _, child := range f.Children {
		err := child.Write(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	rootFile := NewFolder(root, nil, true)
	err := Dirwalk(&rootFile)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(os.Stdout)
	_ = rootFile.Write(w)
	w.Flush()
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
