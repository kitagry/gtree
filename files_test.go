package main

import (
	"bytes"
	"testing"

	"github.com/fatih/color"
)

func TestFileSuffix(t *testing.T) {
	inputs := []struct {
		file   *File
		expect string
	}{
		{NewFile("test.go", nil, false, ""), "go"},
		{NewFile("test.rb", nil, false, ""), "rb"},
		{NewFile("test.html.erb", nil, false, ""), "erb"},
	}

	for _, input := range inputs {
		if input.file.Suffix() != input.expect {
			t.Errorf("file.Suffix() expected %s, but got %s", input.expect, input.file.Suffix())
		}
	}
}

func TestFilesPath(t *testing.T) {
	inputs := []struct {
		file   *File
		expect string
	}{
		{NewFile("test.go", nil, false, ""), "test.go"},
		{NewFile("test.go", NewFolder("test", nil, false), false, ""), "test/test.go"},
		{
			NewFile("test.go",
				NewFolder("test2",
					NewFolder("test1", nil, false),
					false),
				false, ""),
			"test1/test2/test.go",
		},
	}

	for _, input := range inputs {
		if input.file.Path() != input.expect {
			t.Errorf("file.Suffix() expected %s, but got %s", input.expect, input.file.Path())
		}
	}
}

func TestFilesWrite(t *testing.T) {
	goicon, ok := icons["go"]
	if !ok {
		t.Errorf("Go icon is not present")
		return
	}
	goiconString := color.New(goicon.Color).Sprint(goicon.Icon + " ")

	folder := NewFolder("test", nil, false)
	folder.ChildPrefix = "│   "
	inputs := []struct {
		file              *File
		notFullPathExpect string
		fullPathExpect    string
	}{
		{
			NewFile("test.go", nil, false, ""),
			goiconString + "test.go\n",
			goiconString + "test.go\n",
		},
		{
			NewFile("test.go", NewFolder("test", nil, false), false, ""),
			"├── " + goiconString + "test.go\n",
			"├── " + goiconString + "test/test.go\n",
		},
		{
			NewFile("test.go", NewFolder("test", nil, false), true, ""),
			"└── " + goiconString + "test.go\n",
			"└── " + goiconString + "test/test.go\n",
		},
		{
			NewFile("test.go", folder, false, ""),
			"│   ├── " + goiconString + "test.go\n",
			"│   ├── " + goiconString + "test/test.go\n",
		},
		{
			NewFile("test.go", NewFolder("test", nil, false), false, "/test/test/test.go"),
			"├── " + goiconString + symColor.Sprint("test.go") + " -> /test/test/test.go\n",
			"├── " + goiconString + symColor.Sprint("test/test.go") + " -> /test/test/test.go\n",
		},
	}

	for _, in := range inputs {
		buffer := new(bytes.Buffer)
		in.file.Write(buffer, false)
		if buffer.String() != in.notFullPathExpect {
			t.Errorf("file.Write() with not full path expected %s, but got %s", in.notFullPathExpect, buffer.String())
		}

		buffer = new(bytes.Buffer)
		in.file.Write(buffer, true)
		if buffer.String() != in.fullPathExpect {
			t.Errorf("file.Write() with full path expected %s, but got %s", in.fullPathExpect, buffer.String())
		}
	}
}
