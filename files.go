package main

import (
	"errors"
	"io"
	"strings"

	"github.com/fatih/color"
)

var (
	folderColor = color.New(color.FgBlue)
	symColor    = color.New(color.FgHiCyan)
)

// FileInfo is interface for file and folder.
type FileInfo interface {
	// Name returns File name or Folder name.
	Name() string

	// Path returns relative path from file tree's root.
	Path() string

	// IsLast is true when FileInfo is last child of parent's children.
	// This will use to display tree.
	IsLast() bool

	// Write output to io.Writer.
	// If isFullPath is true, Write output fileInfo.Path()
	Write(w io.Writer, isFullPath bool) error
}

type File struct {
	name   string
	parent FileInfo
	isLast bool

	path    string
	SymLink string
}

var _ FileInfo = (*File)(nil)

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
		var err error
		if f.IsLast() {
			_, err = w.Write([]byte(v.ChildPrefix + "└── "))
		} else {
			_, err = w.Write([]byte(v.ChildPrefix + "├── "))
		}

		if err != nil {
			return err
		}

		if f.IsSym() {
			color.New(icon.Color).Fprint(w, icon.Icon+" ")
			var err error
			if isFullPath {
				_, err = symColor.Fprint(w, f.Path())
			} else {
				_, err = symColor.Fprint(w, f.Name())
			}

			if err != nil {
				return err
			}

			_, err = w.Write([]byte(" -> " + f.SymLink + "\n"))
			if err != nil {
				return err
			}
		} else {
			color.New(icon.Color).Fprint(w, icon.Icon+" ")
			var err error
			if isFullPath {
				_, err = w.Write([]byte(f.Path() + "\n"))
			} else {
				_, err = w.Write([]byte(f.Name() + "\n"))
			}

			if err != nil {
				return err
			}
		}
	case nil:
		color.New(icon.Color).Fprint(w, icon.Icon+" ")
		_, err := w.Write([]byte(f.Name() + "\n"))
		if err != nil {
			return err
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
	err  error
}

var _ FileInfo = (*Folder)(nil)

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

func (f *Folder) Path() string {
	if f.path == "" {
		if f.parent == nil {
			f.path = f.Name()
		} else {
			f.path = f.parent.Path() + "/" + f.Name()
		}
	}

	return f.path
}

func (f *Folder) IsLast() bool {
	return f.isLast
}

func (f *Folder) SetError(err error) {
	f.err = err
}

func (f *Folder) displayName() string {
	if f.err != nil {
		return f.Name() + " [" + f.err.Error() + "]"
	}
	return f.Name()
}

func (f *Folder) displayPath() string {
	if f.err != nil {
		return f.Path() + " [" + f.err.Error() + "]"
	}
	return f.Path()
}

func (f *Folder) Write(w io.Writer, isFullPath bool) error {
	if f.parent == nil {
		_, err := folderColor.Fprintln(w, f.Name())
		return err
	}

	switch v := f.parent.(type) {
	case *Folder:
		var err error
		if f.IsLast() {
			_, err = w.Write([]byte(v.ChildPrefix + "└── "))
			f.ChildPrefix = v.ChildPrefix + "    "
		} else {
			_, err = w.Write([]byte(v.ChildPrefix + "├── "))
			f.ChildPrefix = v.ChildPrefix + "│   "
		}

		if err != nil {
			return err
		}

		if isFullPath {
			_, err = folderColor.Fprintln(w, f.displayPath())
		} else {
			_, err = folderColor.Fprintln(w, f.displayName())
		}

		if err != nil {
			return err
		}
	default:
		return errors.New("Unexpected parent type")
	}
	return nil
}
