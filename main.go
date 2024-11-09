package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
)

type RoboHash struct {
	Format      string
	HexDigest   string
	HashArray   []int
	Iter        int
	ResourceDir string
	Sets        []string
	BgSets      []string
	Colors      []string
}

func NewRoboHash(input string, hashCount int, ignoreExt bool) *RoboHash {
	// Create SHA-512
	hash := sha512.New()
	hash.Write([]byte(input))
	hexDigest := hex.EncodeToString(hash.Sum(nil))

	/*
		Start [Iter] at 4, so earlier is reserved
		0 = Color
		1 = Set
		2 = bgset
		3 = BG
	*/

	robohash := &RoboHash{
		Format:    "png",
		HexDigest: hexDigest,
		Iter:      4,
	}

	if ignoreExt {
		robohash.RemoveExts(input)
	}

	robohash.CreateHashes(hashCount)
	robohash.ResourceDir, _ = filepath.Abs(".")

	// Load Directories
	robohash.Sets = robohash.ListDir(filepath.Join(robohash.ResourceDir, "sets"))
	robohash.BgSets = robohash.ListDir(filepath.Join(robohash.ResourceDir, "backgrounds"))
	robohash.Colors = robohash.ListDir(filepath.Join(robohash.ResourceDir, "sets", "set1"))

	return robohash
}

func (r *RoboHash) CreateHashes(count int) {
	/*
		Breaks up our hash into slots, so we can pull them out later.
		Essentially, it splits our SHA/MD5/etc into X parts.
	*/

	blockSize := len(r.HexDigest) / count

	for i := 0; i < count; i++ {
		start := i * blockSize
		end := start + blockSize

		block, _ := hex.DecodeString(r.HexDigest[start:end])
		r.HashArray = append(r.HashArray, int(block[0]))
	}

	// Double the array size
	r.HashArray = append(r.HashArray, r.HashArray...)
}

func (r *RoboHash) ListDir(path string) []string {
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

func (r *RoboHash) GetListOfFiles(path string) []string {
	var chosenFiles []string
	directories := r.ListDir(path)

	for _, dir := range directories {
		dirPath := filepath.Join(path, dir)
		files, _ := filepath.Glob(filepath.Join(dirPath, "*"))
		sort.Strings(files)

		if len(files) > 0 {
			element := r.HashArray[r.Iter] % len(files)
			chosenFiles = append(chosenFiles, files[element])
			r.Iter++
		}
	}

	return chosenFiles
}

func (r *RoboHash) RemoveExts(str string) {
	if ext := filepath.Ext(str); ext != "" {
		str = strings.TrimSuffix(str, ext)
		r.Format = ext[1:]
		if r.Format == "jpg" {
			r.Format = "jpeg"
		}
	}
}

func Overlay(bg, fg image.Image) image.Image {
	output := image.NewRGBA(bg.Bounds())
	draw.Draw(output, output.Bounds(), bg, image.Point{}, draw.Src)
	draw.Draw(output, fg.Bounds(), fg, image.Point{}, draw.Over)
	return output
}

func (r *RoboHash) Assemble(roboset, robocolor, format, bgset string, sizex, sizey int) image.Image {
	/*
	   Build our Robot!
	   Returns the robot image itself.
	*/

	/*
	  Allow users to manually specify a robot 'set' that they like.
	  Ensure that this is one of the allowed choices, or allow all
	  If they don't set one, take the first entry from sets above.
	*/

	if roboset == "any" || roboset == "" {
		roboset = r.Sets[r.HashArray[1]%len(r.Sets)]
	} else {

		found := false

		for _, set := range r.Sets {
			if set == roboset {
				roboset = set
				found = true
				break
			}
		}

		if !found {
			roboset = r.Sets[0]
		}
	}

	/*
	  Only set1 is setup to be color-seletable. The others don't have enough pieces in various colors.
	  This could/should probably be expanded at some point..
	  Right now, this feature is almost never used. ( It was < 44 requests this year, out of 78M reqs )
	*/

	if roboset == "set1" {

		found := false

		for _, color := range r.Colors {
			if robocolor == color {
				roboset = "set1/" + color
				found = true
			}

			if !found {
				randomColor := r.Colors[r.HashArray[0]%len(r.Colors)]
				roboset = "set1/" + randomColor
			}
		}
	}

	// If they specified a background, ensure it's legal, then give it to them.
	if bgset == "any" || bgset == "" {
		bgset = r.BgSets[r.HashArray[2]%len(r.BgSets)]
	} else {

		found := false

		for _, bg := range r.BgSets {
			if bgset == bg {
				bgset = bg
				found = true
			}
		}

		if !found {
			bgset = r.BgSets[r.HashArray[2]%len(r.BgSets)]
		}
	}

	// If we set a format based on extension earlier, use that. Otherwise, PNG.
	if format == "" {
		format = r.Format
	}

	/*
		Each directory in our set represents one piece of the Robot, such as the eyes, nose, mouth, etc.

		Each directory is named with two numbers:
			- The number before the # is the sort order.
		  	This ensures that they always go in the same order when choosing pieces, regardless of OS.

		  - The second number is the order in which to apply the pieces.
		  	For instance, the head has to go down BEFORE the eyes, or the eyes would be hidden.

		First, we'll get a list of parts of our robot.
	*/
	dir := r.ResourceDir + "/sets/" + roboset
	roboparts := r.GetListOfFiles(dir)

	// Now that we've sorted them by the first number, we need to sort each sub-category by the second.
	sort.Slice(roboparts, func(i, j int) bool {
		return strings.Split(roboparts[i], "#")[1] < strings.Split(roboparts[j], "#")[1]
	})

	// *************************

	background := ""
	if bgset != "" {
		path := r.ResourceDir + "/backgrounds/" + bgset
		backgrounds, _ := os.ReadDir(path)

		var bgList []string
		for _, bg := range backgrounds {
			if !strings.HasPrefix(bg.Name(), ".") {
				bgList = append(bgList, path+"/"+bg.Name())
			}
		}
		background = bgList[r.HashArray[3]%len(bgList)]
	}

	// Assemble robot parts
	var roboImg image.Image
	for i, part := range roboparts {
		img, err := imaging.Open(part)
		if err != nil {
			log.Fatal(err)
		}
		img = imaging.Resize(img, 1024, 1024, imaging.Lanczos)

		if i == 0 {
			roboImg = img
		} else {
			roboImg = Overlay(roboImg, img)
		}
	}

	// Apply background if available
	if background != "" {
		bgImg, err := imaging.Open(background)
		if err != nil {
			log.Fatal(err)
		}
		bgImg = imaging.Resize(bgImg, 1024, 1024, imaging.Lanczos)
		roboImg = Overlay(bgImg, roboImg)
	}

	finalImg := imaging.Resize(roboImg, sizex, sizey, imaging.Lanczos)
	return finalImg
}

func main() {
	robo := NewRoboHash("example", 11, true)
	robotImage := robo.Assemble("set1", "blue", "png", "bg2", 300, 300)

	file, err := os.Create("robot.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	png.Encode(file, robotImage)

	fmt.Println("Robot image generated as robot.png")
}

/*
Docs
	roboset = "any", "", "set[1-5]"
	robocolor = only when set1 selected, (blue, brown, green, grey,  orange, pink, purple, red, white, yellow)
	bgset = "any", "",  "bg[1-2]"
*/
