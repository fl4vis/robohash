package utils

import (
	"io/fs"
	"log"
	"sort"
)

func ListDir(fsys fs.FS, path string) []string {
	files, err := fs.ReadDir(fsys, path)
	if err != nil {
		log.Fatal(err)
	}

	var dirs []string

	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f.Name())
		}
	}

	sort.Strings(dirs)
	return dirs
}
