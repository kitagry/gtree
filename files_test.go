package main

import (
	"os"
	"testing"
	"time"
)

type dummyOsFile struct {
	name  string
	isDir bool
	isSym bool
}

func newDummyOsFile(name string, isDir bool, isSym bool) os.FileInfo {
	return &dummyOsFile{
		name:  name,
		isDir: isDir,
		isSym: isSym,
	}
}

func (f *dummyOsFile) Name() string {
	return f.name
}

func (f *dummyOsFile) Size() int64 {
	panic("not implemented")
}

func (f *dummyOsFile) Mode() os.FileMode {
	if f.isSym {
		return os.ModeSymlink
	}
	return os.ModeDir
}

func (f *dummyOsFile) ModTime() time.Time {
	panic("not implemented")
}

func (f *dummyOsFile) IsDir() bool {
	return f.isDir
}

func (f *dummyOsFile) Sys() interface{} {
	panic("not implemented")
}

func TestFileSuffix(t *testing.T) {
	fi1, _ := NewFileInfo(newDummyOsFile("readme", false, false), nil, false)
	fi2, _ := NewFileInfo(newDummyOsFile("test.go", false, false), nil, false)
	fi3, _ := NewFileInfo(newDummyOsFile("test.html.erb", false, false), nil, false)
	inputs := []struct {
		file   FileInfo
		expect string
	}{
		{fi1, "readme"},
		{fi2, "go"},
		{fi3, "erb"},
	}

	for _, input := range inputs {
		if input.file.FileType() != input.expect {
			t.Errorf("file.FileType() expected %s, but got %s", input.expect, input.file.(*file).FileType())
		}
	}
}

func TestFilesPath(t *testing.T) {
	d1, _ := NewFileInfo(newDummyOsFile("test1", true, false), nil, false)
	d2, _ := NewFileInfo(newDummyOsFile("test2", true, false), d1, false)
	// test.go
	fi1, _ := NewFileInfo(newDummyOsFile("test.go", false, false), nil, false)
	// test1 -> test.go
	fi2, _ := NewFileInfo(newDummyOsFile("test.go", false, false), d1, false)
	// test1 -> test2 -> test.go
	fi3, _ := NewFileInfo(newDummyOsFile("test.go", false, false), d2, false)

	inputs := []struct {
		file   FileInfo
		expect string
	}{
		{fi1, "test.go"},
		{fi2, "test1/test.go"},
		{fi3, "test1/test2/test.go"},
	}

	for _, input := range inputs {
		if input.file.Path() != input.expect {
			t.Errorf("file.Path() expected %s, but got %s", input.expect, input.file.Path())
		}
	}
}

func TestFolderPath(t *testing.T) {
	d1, _ := NewFileInfo(newDummyOsFile("test1", true, false), nil, false)
	d2, _ := NewFileInfo(newDummyOsFile("test2", true, false), d1, false)
	// test
	r1, _ := NewFileInfo(newDummyOsFile("test", true, false), nil, false)
	// test1 -> test
	r2, _ := NewFileInfo(newDummyOsFile("test", true, false), d1, false)
	// test1 -> test
	r3, _ := NewFileInfo(newDummyOsFile("test", true, false), d2, false)
	inputs := []struct {
		folder FileInfo
		expect string
	}{
		{r1, "test"},
		{r2, "test1/test"},
		{r3, "test1/test2/test"},
	}

	for _, in := range inputs {
		if in.folder.Path() != in.expect {
			t.Errorf("folder.Path() expected %s, but got %s", in.expect, in.folder.Path())
		}
	}
}
