package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Dirwalk(root FileInfo, ch chan<- FileInfo, listOptions *ListSearchOptions) {
	err := dirwalk(root, ch, 0, listOptions)
	if err != nil {
		fmt.Println(err)
	}
	close(ch)
}

func dirwalk(root FileInfo, ch chan<- FileInfo, depth int, listOptions *ListSearchOptions) error {
	if !root.IsDir() {
		ch <- root
		return nil
	}

	if listOptions.Level != nil && depth >= *listOptions.Level {
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

		child := NewFileInfo(file, root, isLast)

		err = dirwalk(child, ch, depth+1, listOptions)
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
		if inString(f.Name(), opts.IgnorePatterns) {
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

func inString(s string, list []string) bool {
	for _, l := range list {
		if s == l {
			return true
		}
	}
	return false
}
