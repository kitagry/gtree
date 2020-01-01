package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Dirwalk(root FileInfo, ch chan<- FileInfo, listOptions *ListSearchOptions) {
	err := dirwalk(root, ch, listOptions)
	if err != nil {
		fmt.Println(err)
	}
	close(ch)
}

func dirwalk(root FileInfo, ch chan<- FileInfo, listOptions *ListSearchOptions) error {
	switch f := root.(type) {
	case *File:
		ch <- f
	case *Folder:
		files, err := ioutil.ReadDir(f.Path())
		if err != nil {
			// Set error, but display this folder.
			f.SetError(err)
			ch <- f
			return nil
		}
		ch <- f

		files = filterFiles(files, listOptions)

		for i, file := range files {
			isLast := i == len(files)-1

			var child FileInfo
			if file.IsDir() {
				child = NewFolder(file.Name(), f, isLast)
			} else {
				if file.Mode()&os.ModeSymlink != 0 {
					sym, _ := os.Readlink(f.Path() + "/" + file.Name())
					child = NewFile(file.Name(), f, isLast, sym)
				} else {
					child = NewFile(file.Name(), f, isLast, "")
				}
			}

			f.Children = append(f.Children, child)
			err = dirwalk(child, ch, listOptions)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("Unexpected File type")
	}
	return nil
}

// Remove files which don't satisfy options.
func filterFiles(files []os.FileInfo, opts *ListSearchOptions) []os.FileInfo {
	result := make([]os.FileInfo, 0)
	for _, f := range files {
		if opts.IgnorePattern == f.Name() {
			continue
		}

		if opts.IsOnlyDirectry() && !f.IsDir() {
			continue
		}

		if !opts.IsAll() && strings.HasPrefix(f.Name(), ".") {
			continue
		}

		result = append(result, f)
	}
	return result
}
