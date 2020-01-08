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
	if !root.IsDir() {
		ch <- root
		return nil
	}

	files, err := ioutil.ReadDir(root.Path())
	if err != nil {
		root.SetError(err)
		ch <- root
		return nil
	}
	ch <- root

	files = filterFiles(files, listOptions)

	for i, file := range files {
		isLast := i == len(files)-1

		child, err := NewFileInfo(file, root, isLast)
		if err != nil {
			return err
		}

		err = dirwalk(child, ch, listOptions)
		if err != nil {
			return err
		}
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
