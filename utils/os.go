package utils

import (
	"log"
	"os"
	"sort"
)

func ListDir(path string) []string {
	files, err := os.ReadDir(path)
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
