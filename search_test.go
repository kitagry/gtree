package main

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

type dummySearchFileInfo struct {
	name      string
	isDir     bool
	isSymlink bool
}

func newDummySearchFileInfo(name string, isDir, isSymlink bool) *dummySearchFileInfo {
	return &dummySearchFileInfo{
		name:      name,
		isDir:     isDir,
		isSymlink: isSymlink,
	}
}

var _ os.FileInfo = (*dummySearchFileInfo)(nil)

func (d *dummySearchFileInfo) Name() string {
	return d.name
}

func (d *dummySearchFileInfo) Size() int64 {
	return 0
}

func (d *dummySearchFileInfo) Mode() os.FileMode {
	var result os.FileMode
	if d.isDir {
		result |= os.ModeDir
	}

	if d.isSymlink {
		result |= os.ModeSymlink
	}

	return result
}

func (d *dummySearchFileInfo) IsDir() bool {
	return d.isDir
}

func (d *dummySearchFileInfo) Sys() interface{} {
	return nil
}

func (d *dummySearchFileInfo) ModTime() time.Time {
	return time.Now()
}

func (d *dummySearchFileInfo) String() string {
	if d.IsDir() {
		return fmt.Sprintf("Folder(%s)", d.Name())
	}
	return fmt.Sprintf("File(%s)", d.Name())
}

func TestFilterFiles(t *testing.T) {
	files := []os.FileInfo{
		newDummySearchFileInfo("file", false, false),
		newDummySearchFileInfo("sym-file", false, true),
		newDummySearchFileInfo("folder", true, false),
		newDummySearchFileInfo("sym-folder", true, true),
		newDummySearchFileInfo(".dotfile", false, false),
		newDummySearchFileInfo(".sym-dotfile", false, true),
		newDummySearchFileInfo(".dotfolder", true, false),
		newDummySearchFileInfo(".sym-dotfolder", true, true),
	}

	inputs := []struct {
		opts   *ListSearchOptions
		result []os.FileInfo
	}{
		{
			&ListSearchOptions{
				All:           []bool{},
				OnlyDirectory: []bool{},
				IgnorePattern: "",
			},
			[]os.FileInfo{
				newDummySearchFileInfo("file", false, false),
				newDummySearchFileInfo("sym-file", false, true),
				newDummySearchFileInfo("folder", true, false),
				newDummySearchFileInfo("sym-folder", true, true),
			},
		},
		{
			&ListSearchOptions{
				All:           []bool{true},
				OnlyDirectory: []bool{},
				IgnorePattern: "",
			},
			[]os.FileInfo{
				newDummySearchFileInfo("file", false, false),
				newDummySearchFileInfo("sym-file", false, true),
				newDummySearchFileInfo("folder", true, false),
				newDummySearchFileInfo("sym-folder", true, true),
				newDummySearchFileInfo(".dotfile", false, false),
				newDummySearchFileInfo(".sym-dotfile", false, true),
				newDummySearchFileInfo(".dotfolder", true, false),
				newDummySearchFileInfo(".sym-dotfolder", true, true),
			},
		},
		{
			&ListSearchOptions{
				All:           []bool{},
				OnlyDirectory: []bool{true},
				IgnorePattern: "",
			},
			[]os.FileInfo{
				newDummySearchFileInfo("folder", true, false),
				newDummySearchFileInfo("sym-folder", true, true),
			},
		},
		{
			&ListSearchOptions{
				All:           []bool{true},
				OnlyDirectory: []bool{true},
				IgnorePattern: "",
			},
			[]os.FileInfo{
				newDummySearchFileInfo("folder", true, false),
				newDummySearchFileInfo("sym-folder", true, true),
				newDummySearchFileInfo(".dotfolder", true, false),
				newDummySearchFileInfo(".sym-dotfolder", true, true),
			},
		},
		{
			&ListSearchOptions{
				All:           []bool{true},
				OnlyDirectory: []bool{true},
				IgnorePattern: "folder",
			},
			[]os.FileInfo{
				newDummySearchFileInfo("sym-folder", true, true),
				newDummySearchFileInfo(".dotfolder", true, false),
				newDummySearchFileInfo(".sym-dotfolder", true, true),
			},
		},
	}

	for i, in := range inputs {
		res := filterFiles(files, in.opts)
		if !reflect.DeepEqual(res, in.result) {
			t.Errorf("%d: filterFiles number expected %v, got %v", i, in.result, res)
		}
	}

}
