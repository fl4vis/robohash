package generator

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/fs"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

const workingImageSize = 1024

func overlay(background, foreground image.Image) image.Image {
	output := image.NewRGBA(background.Bounds())
	draw.Draw(output, output.Bounds(), background, image.Point{}, draw.Src)
	draw.Draw(output, output.Bounds(), foreground, image.Point{}, draw.Over)
	return output
}

func decodeImage(fsys fs.FS, path string) (image.Image, error) {
	file, err := fsys.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open image %q: %w", path, err)
	}
	defer file.Close()

	img, err := imaging.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode image %q: %w", path, err)
	}
	return img, nil
}

func resizeWorkingImage(img image.Image) image.Image {
	return imaging.Resize(img, workingImageSize, workingImageSize, imaging.Lanczos)
}

func prepareOutput(img image.Image, format string, width, height int) image.Image {
	if format == "jpeg" || format == "jpg" {
		bounds := img.Bounds()
		flattened := image.NewRGBA(bounds)
		draw.Draw(flattened, bounds, &image.Uniform{C: color.White}, image.Point{}, draw.Src)
		draw.Draw(flattened, bounds, img, image.Point{}, draw.Over)
		img = flattened
	}
	return imaging.Resize(img, width, height, imaging.Lanczos)
}
