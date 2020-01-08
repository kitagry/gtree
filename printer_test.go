package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
)

type dummyPrinterFileInfo struct {
	name        string
	path        string
	filetype    string
	childPrefix string
	symlink     string
	isLast      bool
	isDir       bool
	parent      FileInfo
	err         error
}

func newDummyPrinterFileInfo(name, path, filetype, childPrefix, symlink string, isLast, isDir bool, err error, parent FileInfo) FileInfo {
	return &dummyPrinterFileInfo{
		name:        name,
		path:        path,
		filetype:    filetype,
		childPrefix: childPrefix,
		symlink:     symlink,
		isLast:      isLast,
		isDir:       isDir,
		parent:      parent,
		err:         err,
	}
}

func (d *dummyPrinterFileInfo) Name() string {
	return d.name
}

func (d *dummyPrinterFileInfo) Path() string {
	return d.path
}

func (d *dummyPrinterFileInfo) FileType() string {
	return d.filetype
}

func (d *dummyPrinterFileInfo) IsLast() bool {
	return d.isLast
}

func (d *dummyPrinterFileInfo) IsSym() bool {
	return d.symlink != ""
}

func (d *dummyPrinterFileInfo) SymLink() (string, error) {
	if d.symlink == "" {
		return "", errors.New("This is not SymLink")
	}
	return d.symlink, nil
}

func (d *dummyPrinterFileInfo) Parent() (FileInfo, bool) {
	if d.parent == nil {
		return nil, false
	}
	return d.parent, true
}

func (d *dummyPrinterFileInfo) ChildPrefix() string {
	return d.childPrefix
}

func (d *dummyPrinterFileInfo) Write(w io.Writer, isFullPath bool) error {
	return nil
}

func (d *dummyPrinterFileInfo) IsDir() bool {
	return d.isDir
}

func (d *dummyPrinterFileInfo) SetError(err error) {
}

func (d *dummyPrinterFileInfo) Error() error {
	return d.err
}

func TestPrinter_Write(t *testing.T) {
	p := NewPrinter()
	inputs := []struct {
		fileInfo            FileInfo
		fullPathExpected    string
		notFullPathExpected string
	}{
		{
			newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", false, false, nil, nil),
			NewIconString("go") + " test.go\n",
			NewIconString("go") + " test/test.go\n",
		},
		{
			newDummyPrinterFileInfo("test.rb", "test/test.rb", "rb", "", "", false, false, nil, nil),
			NewIconString("rb") + " test.rb\n",
			NewIconString("rb") + " test/test.rb\n",
		},
		{
			newDummyPrinterFileInfo("test", "test/test", "", "", "", false, true, nil, nil),
			folderColor.Sprintf("test") + "\n",
			folderColor.Sprintf("test/test") + "\n",
		},
		{
			newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", false, false, nil,
				newDummyPrinterFileInfo("test", "test", "", "", "", false, true, nil, nil),
			),
			"├── " + NewIconString("go") + " test.go\n",
			"├── " + NewIconString("go") + " test/test.go\n",
		},
		{
			newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", true, false, nil,
				newDummyPrinterFileInfo("test", "test", "", "", "", false, true, nil, nil),
			),
			"└── " + NewIconString("go") + " test.go\n",
			"└── " + NewIconString("go") + " test/test.go\n",
		},
		{
			newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", true, false, nil,
				newDummyPrinterFileInfo("test", "test", "", "│   ", "", false, true, nil, nil),
			),
			"│   └── " + NewIconString("go") + " test.go\n",
			"│   └── " + NewIconString("go") + " test/test.go\n",
		},
		{
			newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "/test.go", false, false, nil, nil),
			fmt.Sprintf("%s %s -> %s\n", NewIconString("go"), symColor.Sprint("test.go"), "/test.go"),
			fmt.Sprintf("%s %s -> %s\n", NewIconString("go"), symColor.Sprint("test/test.go"), "/test.go"),
		},
		{
			newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", false, false, errors.New("Hello"), nil),
			NewIconString("go") + " test.go [Hello]\n",
			NewIconString("go") + " test/test.go [Hello]\n",
		},
	}

	for _, input := range inputs {
		buffer := new(bytes.Buffer)
		p.Write(buffer, input.fileInfo, false)
		if buffer.String() != input.fullPathExpected {
			t.Errorf("printer.Write() expected %s, got %s", input.fullPathExpected, buffer.String())
		}

		buffer = new(bytes.Buffer)
		p.Write(buffer, input.fileInfo, true)
		if buffer.String() != input.notFullPathExpected {
			t.Errorf("printer.Write() expected %s, got %s", input.notFullPathExpected, buffer.String())
		}
	}
}
