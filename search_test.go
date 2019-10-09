package main

import (
	"os"
	"reflect"
	"testing"
	"time"
)

type dummyFileInfo struct {
	name      string
	isDir     bool
	isSymlink bool
}

func newDummyFileInfo(name string, isDir, isSymlink bool) *dummyFileInfo {
	return &dummyFileInfo{
		name:      name,
		isDir:     isDir,
		isSymlink: isSymlink,
	}
}

var _ os.FileInfo = (*dummyFileInfo)(nil)

func (d *dummyFileInfo) Name() string {
	return d.name
}

func (d *dummyFileInfo) Size() int64 {
	return 0
}

func (d *dummyFileInfo) Mode() os.FileMode {
	var result os.FileMode
	if d.isDir {
		result |= os.ModeDir
	}

	if d.isSymlink {
		result |= os.ModeSymlink
	}

	return result
}

func (d *dummyFileInfo) IsDir() bool {
	return d.isDir
}

func (d *dummyFileInfo) Sys() interface{} {
	return nil
}

func (d *dummyFileInfo) ModTime() time.Time {
	return time.Now()
}

func TestFilterFiles(t *testing.T) {
	files := []os.FileInfo{
		newDummyFileInfo("file", false, false),
		newDummyFileInfo("sym-file", false, true),
		newDummyFileInfo("folder", true, false),
		newDummyFileInfo("sym-folder", true, true),
		newDummyFileInfo(".dotfile", false, false),
		newDummyFileInfo(".sym-dotfile", false, true),
		newDummyFileInfo(".dotfolder", true, false),
		newDummyFileInfo(".sym-dotfolder", true, true),
	}

	inputs := []struct {
		opts   *ListSearchOptions
		result []os.FileInfo
	}{
		{
			&ListSearchOptions{
				All:           []bool{},
				OnlyDirectory: []bool{},
			},
			[]os.FileInfo{
				newDummyFileInfo("file", false, false),
				newDummyFileInfo("sym-file", false, true),
				newDummyFileInfo("folder", true, false),
				newDummyFileInfo("sym-folder", true, true),
			},
		},
		{
			&ListSearchOptions{
				All:           []bool{true},
				OnlyDirectory: []bool{},
			},
			[]os.FileInfo{
				newDummyFileInfo("file", false, false),
				newDummyFileInfo("sym-file", false, true),
				newDummyFileInfo("folder", true, false),
				newDummyFileInfo("sym-folder", true, true),
				newDummyFileInfo(".dotfile", false, false),
				newDummyFileInfo(".sym-dotfile", false, true),
				newDummyFileInfo(".dotfolder", true, false),
				newDummyFileInfo(".sym-dotfolder", true, true),
			},
		},
		{
			&ListSearchOptions{
				All:           []bool{},
				OnlyDirectory: []bool{true},
			},
			[]os.FileInfo{
				newDummyFileInfo("folder", true, false),
				newDummyFileInfo("sym-folder", true, true),
			},
		},
		{
			&ListSearchOptions{
				All:           []bool{true},
				OnlyDirectory: []bool{true},
			},
			[]os.FileInfo{
				newDummyFileInfo("folder", true, false),
				newDummyFileInfo("sym-folder", true, true),
				newDummyFileInfo(".dotfolder", true, false),
				newDummyFileInfo(".sym-dotfolder", true, true),
			},
		},
	}

	for _, in := range inputs {
		res := filterFiles(files, in.opts)
		if !reflect.DeepEqual(res, in.result) {
			t.Errorf("filterFiles number expected %d, got %d", len(in.result), len(res))
		}
	}

}
