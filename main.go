package main

import (
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/fl4vis/robohash/utils"
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
	Img         image.Image
}

func NewRoboHash(input string, hashCount int, ignoreExt bool) *RoboHash {
	/*
		Start [Iter] at 4, so earlier is reserved
		0 = Color
		1 = Set
		2 = bgset
		3 = BG
	*/

	robohash := &RoboHash{
		Format: "png",
		Iter:   4,
	}

	// Optionally remove an images extension before hashing.
	if ignoreExt {
		input = robohash.RemoveExts(input)
	}

	// Create SHA-512
	hash := sha512.New()
	hash.Write([]byte(input))
	hexDigest := hex.EncodeToString(hash.Sum(nil))

	robohash.HexDigest = hexDigest

	robohash.CreateHashes(hashCount)
	robohash.ResourceDir, _ = filepath.Abs("./img/robo/")

	// Load Directories
	robohash.Sets = utils.ListDir(filepath.Join(robohash.ResourceDir, "sets"))
	robohash.BgSets = utils.ListDir(filepath.Join(robohash.ResourceDir, "backgrounds"))
	robohash.Colors = utils.ListDir(filepath.Join(robohash.ResourceDir, "sets", "set1"))

	return robohash
}

func (r *RoboHash) CreateHashes(count int) {
	/*
		Breaks up our hash into slots, so we can pull them out later.
		Essentially, it splits our SHA/MD5/etc into X parts.
	*/

	// Get 1/numblocks of the hash
	blockSize := len(r.HexDigest) / count

	for i := 0; i < count; i++ {
		start := i * blockSize
		end := start + blockSize

		block, _ := hex.DecodeString(r.HexDigest[start:end])
		r.HashArray = append(r.HashArray, int(block[0]))
	}

	// Double the array size
	/*
	  Workaround for adding more sets in 2019.
	  We run out of blocks, because we use some for each set, whether it's called or not.
	  I can't easily change this without invalidating every hash so far :/
	  This shouldn't reduce the security since it should only draw from one set of these in practice.
	*/
	r.HashArray = append(r.HashArray, r.HashArray...)
}

func (r *RoboHash) GetListOfFiles(path string) []string {
	/*
	  Go through each subdirectory of `path`, and choose one file from each to use in our hash.
	  Continue to increase self.iter, so we use a different 'slot' of randomness each time.
	*/
	var chosenFiles []string

	// Get a list of all subdirectories
	directories := utils.ListDir(path)

	// Go through each directory in the list, and choose one file from each.
	// Add this file to our master list of robotparts.
	for _, dir := range directories {
		dirPath := filepath.Join(path, dir)
		files, _ := filepath.Glob(filepath.Join(dirPath, "*"))
		sort.Strings(files)

		if len(files) > 0 {
			// Use some of our hash bits to choose which file
			element := r.HashArray[r.Iter] % len(files)
			chosenFiles = append(chosenFiles, files[element])
			r.Iter++
		}
	}

	return chosenFiles
}

func (r *RoboHash) RemoveExts(str string) string {
	// Sets the string, to create the Robohash

	/*
	  If the user hasn't disabled it, we will detect image extensions, such as .png, .jpg, etc.
	  We'll remove them from the string before hashing.
	  This ensures that /Bear.png and /Bear.bmp will send back the same image, in different formats.
	*/

	str = strings.ToLower(str)
	suffixes := []string{".png", ".gif", ".jpg", ".bmp", ".jpeg", ".ppm", ".datauri"}

	for _, suffix := range suffixes {
		if strings.HasSuffix(str, suffix) {
			startIndex := strings.LastIndex(str, suffix)
			format := str[startIndex:]

			if format == ".jpg" {
				format = ".jpeg"
			}

			str = str[0:startIndex] + format
		}
	}

	return str
}

func Overlay(bg, fg image.Image) image.Image {
	output := image.NewRGBA(bg.Bounds())
	draw.Draw(output, output.Bounds(), bg, image.Point{}, draw.Src)
	draw.Draw(output, fg.Bounds(), fg, image.Point{}, draw.Over)
	return output
}

func (r *RoboHash) Assemble(roboset, robocolor, format, bgset string, sizex, sizey int) *RoboHash {
	/*
	   Build our Robot!
	   Returns the robot image itself.
	*/

	/*
	  Allow users to manually specify a robot 'set' that they like.
	  Ensure that this is one of the allowed choices, or allow all
	  If they don't set one, take the first entry from sets above.
	*/

	if roboset == "any" {
		roboset = r.Sets[r.HashArray[1]%len(r.Sets)]
	} else {

		found := false
		roboset = "set" + roboset

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
	switch bgset {
	case "any":
		bgset = r.BgSets[r.HashArray[2]%len(r.BgSets)]
	case "":
		bgset = ""
	default:

		found := false
		bgset = "bg" + bgset

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

	// Flatten if JPEG
	if format == "jpeg" {
		rect := roboImg.Bounds()
		flattenedImg := image.NewRGBA(rect)
		draw.Draw(flattenedImg, rect, &image.Uniform{color.White}, image.Point{}, draw.Src)
		draw.Draw(flattenedImg, rect, roboImg, image.Point{}, draw.Over)
		roboImg = flattenedImg
	}

	// Resize final image formatand assign it to r.Img
	r.Img = imaging.Resize(roboImg, sizex, sizey, imaging.Lanczos)
	r.Format = format

	return r
}

func main() {
	extension := utils.ExtensionsType{
		ValidValues: []string{"png", "jpg", "jpeg", "gif", "ppm", "datauri"},
		Value:       "png",
	}

	// Flags
	input := flag.String("input", "", "String to be hashed")
	hashSlots := flag.Int("slots", 11, "Number of slots hash will be divided")
	ignoreExt := flag.Bool("ignoreExt", true, "Removes any extension of the input before hashing")

	set := flag.String("set", "any", "[1,5] || any")
	color := flag.String("color", "", "Only when set1 selected -> [blue, brown, green, grey,  orange, pink, purple, red, white, yellow]")
	flag.Var(&extension, "extension", "[png, jpg, jpeg, gif, ppm, datauri]")
	bgset := flag.String("bgset", "", "[1,2] || any")
	sizex := flag.Int("sizex", 300, "size x")
	sizey := flag.Int("sizey", 300, "size y")

	flag.Parse()

	log.SetFlags(0)
	if *input == "" {
		log.Fatalf("Empty input value not allowed")
	}

	if *hashSlots > 30 {
		log.Fatalf("Maxium number of recommended slots reached: %d", *hashSlots)
	}

	robo := NewRoboHash(*input, *hashSlots, *ignoreExt)
	robotImage := robo.Assemble(*set, *color, extension.Value, *bgset, *sizex, *sizey)

	outputFile := fmt.Sprintf("%s.%s", *input, robo.Format)
	if err := utils.SaveImage(outputFile, robo.Format, robotImage.Img); err != nil {
		log.Fatal(err)
	}
}

/*
Docs
	roboset = "any", "", "set[1-5]"
	robocolor = only when set1 selected, (blue, brown, green, grey,  orange, pink, purple, red, white, yellow)
	bgset = "any", "",  "bg[1-2]"
*/
