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
	tests := map[string]struct {
		fileInfo      FileInfo
		displayOption *ListDisplayOptions
		output        string
	}{
		"print normal": {
			fileInfo: newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", false, false, nil, nil),
			displayOption: &ListDisplayOptions{
				FullPath: nil,
				NoIcons:  nil,
			},
			output: NewIconString("go") + " test.go\n",
		},
		"print full path": {
			fileInfo: newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", false, false, nil, nil),
			displayOption: &ListDisplayOptions{
				FullPath: []bool{true},
				NoIcons:  nil,
			},
			output: NewIconString("go") + " test/test.go\n",
		},
		"print directory": {
			fileInfo: newDummyPrinterFileInfo("test", "test/test", "", "", "", false, true, nil, nil),
			displayOption: &ListDisplayOptions{
				FullPath: nil,
				NoIcons:  nil,
			},
			output: folderColor.Sprintf("%s %s", defaultFolderIcon.Icon, "test") + "\n",
		},
		"print child file": {
			fileInfo: newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", false, false, nil,
				newDummyPrinterFileInfo("test", "test", "", "", "", false, true, nil, nil),
			),
			displayOption: &ListDisplayOptions{
				FullPath: nil,
				NoIcons:  nil,
			},
			output: "├── " + NewIconString("go") + " test.go\n",
		},
		"print last child file": {
			fileInfo: newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", true, false, nil,
				newDummyPrinterFileInfo("test", "test", "", "", "", false, true, nil, nil),
			),
			displayOption: &ListDisplayOptions{
				FullPath: nil,
				NoIcons:  nil,
			},
			output: "└── " + NewIconString("go") + " test.go\n",
		},
		"print grandchild file": {
			fileInfo: newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", true, false, nil,
				newDummyPrinterFileInfo("test", "test", "", "│   ", "", false, true, nil, nil),
			),
			displayOption: &ListDisplayOptions{
				FullPath: nil,
				NoIcons:  nil,
			},
			output: "│   └── " + NewIconString("go") + " test.go\n",
		},
		"print symlink file": {
			fileInfo: newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "/test.go", false, false, nil, nil),
			displayOption: &ListDisplayOptions{
				FullPath: nil,
				NoIcons:  nil,
			},
			output: fmt.Sprintf("%s %s -> %s\n", NewIconString("go"), symColor.Sprint("test.go"), "/test.go"),
		},
		"print error file": {
			fileInfo: newDummyPrinterFileInfo("test.go", "test/test.go", "go", "", "", false, false, errors.New("Hello"), nil),
			displayOption: &ListDisplayOptions{
				FullPath: nil,
				NoIcons:  nil,
			},
			output: NewIconString("go") + " test.go [Hello]\n",
		},
	}

	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			p := NewPrinter(tt.displayOption)

			buffer := new(bytes.Buffer)
			p.Write(buffer, tt.fileInfo)
			if buffer.String() != tt.output {
				t.Errorf("printer.Write() expected '%s', got '%s'", tt.output, buffer.String())
			}
		})
	}
}
