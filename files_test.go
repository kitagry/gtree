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

func newDummyOsFile(name string, isDir bool) os.FileInfo {
	return &dummyOsFile{
		name:  name,
		isDir: isDir,
		isSym: false,
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

func TestFileInfo_FileType(t *testing.T) {
	fi1 := NewFileInfo(newDummyOsFile("readme", false), nil, false)
	fi2 := NewFileInfo(newDummyOsFile("test.go", false), nil, false)
	fi3 := NewFileInfo(newDummyOsFile("test.html.erb", false), nil, false)
	fo := NewFileInfo(newDummyOsFile("directory", true), nil, false)
	inputs := []struct {
		file   FileInfo
		expect string
	}{
		{fi1, "readme"},
		{fi2, "go"},
		{fi3, "erb"}, // TODO: Change "html.erb"?
		{fo, ""},
	}

	for _, input := range inputs {
		if input.file.FileType() != input.expect {
			t.Errorf("file.FileType() expected %s, but got %s", input.expect, input.file.(*file).FileType())
		}
	}
}

func TestFileInfo_Path(t *testing.T) {
	d1 := NewFileInfo(newDummyOsFile("test1", true), nil, false)
	d2 := NewFileInfo(newDummyOsFile("test2", true), d1, false)
	// test.go
	fi1 := NewFileInfo(newDummyOsFile("test.go", false), nil, false)
	// test1 -> test.go
	fi2 := NewFileInfo(newDummyOsFile("test.go", false), d1, false)
	// test1 -> test2 -> test.go
	fi3 := NewFileInfo(newDummyOsFile("test.go", false), d2, false)

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

func TestFileInfo_Parent(t *testing.T) {
	d1 := NewFileInfo(newDummyOsFile("test1", true), nil, false)
	f1 := NewFileInfo(newDummyOsFile("test.go", false), nil, false)
	f2 := NewFileInfo(newDummyOsFile("test.go", false), d1, false)

	inputs := []struct {
		file   FileInfo
		parent FileInfo
		ok     bool
	}{
		{f1, nil, false},
		{f2, d1, true},
	}

	for _, input := range inputs {
		p, ok := input.file.Parent()
		if p != input.parent {
			t.Errorf("FileInfo %v parent expect %v, but got %v", input.file, input.parent, p)
		}

		if ok != input.ok {
			t.Errorf("FileInfo %v Parent() ok expect %v, but got %v", input.file, input.ok, ok)
		}
	}
}

func TestFileInfo_ChildPrefix(t *testing.T) {
	root := NewFileInfo(newDummyOsFile("root", true), nil, false)
	f1 := NewFileInfo(newDummyOsFile("test.go", false), nil, false)
	f2 := NewFileInfo(newDummyOsFile("test.go", false), root, false)
	d1 := NewFileInfo(newDummyOsFile("test", true), root, false)
	d2 := NewFileInfo(newDummyOsFile("test", true), root, true)

	inputs := []struct {
		file        FileInfo
		childPrefix string
	}{
		{root, ""},
		{f1, ""},
		{f2, ""},
		{d1, "│   "},
		{d2, "    "},
	}

	for _, input := range inputs {
		if input.file.ChildPrefix() != input.childPrefix {
			t.Errorf("FileInfo.ChildPrefix expect %s, but got %s", input.childPrefix, input.file.ChildPrefix())
		}
	}
}
