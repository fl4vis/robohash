package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/fl4vis/robohash/utils"
)

type GhibliHash struct {
	HexDigest   string
	HashArray   []int
	Iter        int
	ResourceDir string
	Sets        []string
	BgSets      []string
	Colors      []string
}

func NewGhibliHash() *GhibliHash {
	giblihash := &GhibliHash{}

	giblihash.ResourceDir, _ = filepath.Abs("./img/ghibli")

	return giblihash
}

func main() {
	g := NewGhibliHash()

	dirs := utils.ListDir(g.ResourceDir)

	result := g.GetListOfFiles(dirs)
	fmt.Println(result)

	file, _ := os.Open(result)

	defer file.Close()

	image, _, _ := image.Decode(file)
	// resize := imaging.Resize(file, 1024, 1024, imaging.Lanczos)
	// utils.SaveImage("hi.png", "png", resize)
	utils.SaveImage("hi.png", "png", image)
}

func (g *GhibliHash) GetListOfFiles(directories []string) string {
	dirPath := filepath.Join(g.ResourceDir, directories[0])
	matches, err := filepath.Glob(filepath.Join(dirPath, "*"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(matches)
	return matches[0]
}
