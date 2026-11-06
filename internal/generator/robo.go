package generator

import (
	"fmt"
	"image"
	"io/fs"
	"sort"
	"strings"
)

const roboRoot = "img/robo"

// RoboOptions controls the composition of a RoboHash avatar.
type RoboOptions struct {
	Set        string
	Color      string
	Background string
	Format     string
	Width      int
	Height     int
}

// GenerateRobo builds a deterministic RoboHash avatar from the supplied assets.
func GenerateRobo(fsys fs.FS, input string, slots int, ignoreExtension bool, options RoboOptions) (image.Image, error) {
	hash, err := NewHashSlots(input, slots, ignoreExtension)
	if err != nil {
		return nil, err
	}

	sets, err := directories(fsys, roboRoot+"/sets")
	if err != nil {
		return nil, err
	}
	if len(sets) == 0 {
		return nil, fmt.Errorf("no RoboHash sets found")
	}
	backgroundSets, err := directories(fsys, roboRoot+"/backgrounds")
	if err != nil {
		return nil, err
	}
	colors, err := directories(fsys, roboRoot+"/sets/set1")
	if err != nil {
		return nil, err
	}

	set := selectSet(options.Set, sets, hash)
	if set == "set1" {
		set = "set1/" + selectColor(options.Color, colors, hash)
	}
	backgroundSet := selectBackground(options.Background, backgroundSets, hash)

	parts, err := selectedPartFiles(fsys, roboRoot+"/sets/"+set, hash)
	if err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, fmt.Errorf("no RoboHash parts found for %q", set)
	}

	var result image.Image
	for index, path := range parts {
		part, err := decodeImage(fsys, path)
		if err != nil {
			return nil, err
		}
		part = resizeWorkingImage(part)
		if index == 0 {
			result = part
		} else {
			result = overlay(result, part)
		}
	}

	if backgroundSet != "" {
		background, err := selectedFile(fsys, roboRoot+"/backgrounds/"+backgroundSet, hash.At(3))
		if err != nil {
			return nil, err
		}
		backgroundImage, err := decodeImage(fsys, background)
		if err != nil {
			return nil, err
		}
		result = overlay(resizeWorkingImage(backgroundImage), result)
	}

	return prepareOutput(result, options.Format, options.Width, options.Height), nil
}

func directories(fsys fs.FS, path string) ([]string, error) {
	entries, err := fs.ReadDir(fsys, path)
	if err != nil {
		return nil, fmt.Errorf("read directory %q: %w", path, err)
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

func selectedPartFiles(fsys fs.FS, path string, hash *HashSlots) ([]string, error) {
	directories, err := directories(fsys, path)
	if err != nil {
		return nil, err
	}

	paths := make([]string, 0, len(directories))
	for index, directory := range directories {
		path, err := selectedFile(fsys, path+"/"+directory, hash.At(index+4))
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	sort.Slice(paths, func(i, j int) bool {
		return layerOrder(paths[i]) < layerOrder(paths[j])
	})
	return paths, nil
}

func selectedFile(fsys fs.FS, path string, value int) (string, error) {
	entries, err := fs.ReadDir(fsys, path)
	if err != nil {
		return "", fmt.Errorf("read directory %q: %w", path, err)
	}
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			files = append(files, path+"/"+entry.Name())
		}
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no image files found in %q", path)
	}
	sort.Strings(files)
	return files[value%len(files)], nil
}

func selectSet(requested string, sets []string, hash *HashSlots) string {
	if requested == "any" {
		return sets[hash.At(1)%len(sets)]
	}
	requested = "set" + requested
	for _, set := range sets {
		if set == requested {
			return set
		}
	}
	return sets[0]
}

func selectColor(requested string, colors []string, hash *HashSlots) string {
	for _, color := range colors {
		if color == requested {
			return color
		}
	}
	return colors[hash.At(0)%len(colors)]
}

func selectBackground(requested string, backgrounds []string, hash *HashSlots) string {
	switch requested {
	case "":
		return ""
	case "any":
		return backgrounds[hash.At(2)%len(backgrounds)]
	}
	requested = "bg" + requested
	for _, background := range backgrounds {
		if background == requested {
			return background
		}
	}
	return backgrounds[hash.At(2)%len(backgrounds)]
}

func layerOrder(path string) string {
	parts := strings.Split(path, "#")
	if len(parts) < 2 {
		return path
	}
	return parts[1]
}
