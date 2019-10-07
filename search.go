package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Dirwalk(root FileInfo, ch chan<- FileInfo, listOptions *ListSearchOptions) {
	dirwalk(root, ch, listOptions)
	close(ch)
}

func dirwalk(root FileInfo, ch chan<- FileInfo, listOptions *ListSearchOptions) {
	ch <- root
	switch f := root.(type) {
	case *File:
		return
	case *Folder:
		files, err := ioutil.ReadDir(f.Path())
		if err != nil {
			fmt.Println(err)
			close(ch)
			return
		}

		if len(listOptions.OnlyDirectory) != 0 {
			tmp := make([]os.FileInfo, 0)
			for _, f := range files {
				if f.IsDir() {
					tmp = append(tmp, f)
				}
			}
			files = tmp
		}

		for i, file := range files {
			if len(listOptions.IsAll) == 0 && strings.HasPrefix(file.Name(), ".") {
				continue
			}

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
			dirwalk(child, ch, listOptions)
		}
	default:
		fmt.Println("Unexpected File type")
		close(ch)
		return
	}
}
