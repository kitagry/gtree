package main

import (
	"os"
	"strings"

	"golang.org/x/xerrors"
)

// FileInfo is interface for file and folder.
type FileInfo interface {
	// Name returns File name or Folder name.
	Name() string

	// Path returns relative path from file tree's root.
	Path() string

	// FileType returns files type
	// If FileInfo is directory, returns ""
	FileType() string

	// IsLast is true when FileInfo is last child of parent's children.
	// This will use to display tree.
	IsLast() bool

	// Parent returns fileinfo's parent.
	// When FileInfo doesn't have parent, return false.
	Parent() (FileInfo, bool)

	// ChildPrefix is prefix of children's prefix
	ChildPrefix() string

	// IsDir returns true, when FileInfo is directory
	IsDir() bool

	// IsSym returns true, when FileInfo is symlink
	IsSym() bool

	// SymLink returns symlink path
	SymLink() (string, error)

	// SetError set error
	SetError(err error)

	// Error return error
	Error() error
}

func NewFileInfo(f os.FileInfo, parent FileInfo, isLast bool) (FileInfo, error) {
	var result FileInfo
	if f.IsDir() {
		result = newFolder(f, parent, isLast)
	} else {
		result = newFile(f, parent, isLast)
	}
	return result, nil
}

type baseFileInfo struct {
	os.FileInfo

	parent FileInfo
	isLast bool
	path   string
	err    error
}

func (f *baseFileInfo) Path() string {
	if f.path == "" {
		if f.parent == nil {
			f.path = f.Name()
		} else {
			f.path = f.parent.Path() + "/" + f.Name()
		}
	}

	return f.path
}

func (f *baseFileInfo) Parent() (FileInfo, bool) {
	if f.parent == nil {
		return nil, false
	}
	return f.parent, true
}

func (f *baseFileInfo) IsLast() bool {
	return f.isLast
}

func (f *baseFileInfo) IsSym() bool {
	return f.Mode()&os.ModeSymlink != 0
}

func (f *baseFileInfo) SymLink() (string, error) {
	if !f.IsSym() {
		// TODO: it may returns error
		return "", xerrors.New("This is not symlink")
	}

	symLink, err := os.Readlink(f.Path())
	if err != nil {
		return "", err
	}

	return symLink, nil
}

func (f *baseFileInfo) SetError(err error) {
	f.err = err
}

func (f *baseFileInfo) Error() error {
	return f.err
}

type file struct {
	baseFileInfo
}

var _ FileInfo = (*file)(nil)

func newFile(f os.FileInfo, parent FileInfo, isLast bool) FileInfo {
	return &file{
		baseFileInfo{
			FileInfo: f,
			parent:   parent,
			isLast:   isLast,
		},
	}
}

func (f *file) FileType() string {
	n := strings.Split(f.Name(), ".")
	return n[len(n)-1]
}

func (f *file) ChildPrefix() string {
	return ""
}

func (f *file) IsSym() bool {
	return f.Mode()&os.ModeSymlink != 0
}

func (f *file) symLink() (string, error) {
	symLink, err := os.Readlink(f.Path())
	if err != nil {
		return "", err
	}
	return symLink, nil
}

type folder struct {
	baseFileInfo

	childPrefix string
}

var _ FileInfo = (*folder)(nil)

func newFolder(f os.FileInfo, parent FileInfo, isLast bool) FileInfo {
	return &folder{
		baseFileInfo: baseFileInfo{
			FileInfo: f,
			parent:   parent,
			isLast:   isLast,
		},
	}
}

func (f *folder) FileType() string {
	return ""
}

func (f *folder) ChildPrefix() string {
	p, ok := f.Parent()
	if f.childPrefix == "" && ok {
		if f.IsLast() {
			f.childPrefix = p.ChildPrefix() + "    "
		} else {
			f.childPrefix = p.ChildPrefix() + "│   "
		}
	}
	return f.childPrefix
}
