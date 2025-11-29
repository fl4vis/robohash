package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"

	_ "golang.org/x/image/webp"
)

// saveImage saves an image.Image in the specified format
func SaveImage(filename, format string, img image.Image) error {
	if format == "datauri" {
		dataURI, err := saveDataURI(img)
		if err != nil {
			return err
		}
		fmt.Println(dataURI)
		return nil
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "png":
		return png.Encode(file, img)
	case "jpeg", "jpg":
		options := &jpeg.Options{Quality: 90}
		return jpeg.Encode(file, img, options)
	case "gif":
		return gif.Encode(file, img, nil)
	case "ppm":
		return encodePPM(file, img)
	case "webp":
		return encodeWep(file, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// encodePPM encodes an image in the PPM (Portable Pixmap) format.
func encodePPM(w *os.File, img image.Image) error {
	bounds := img.Bounds()
	_, err := fmt.Fprintf(w, "P6\n%d %d\n255\n", bounds.Dx(), bounds.Dy())
	if err != nil {
		return err
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if _, err := w.Write([]byte{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}); err != nil {
				return err
			}
		}
	}
	return nil
}

// saveDataURI encodes the image to a base64 Data URI and writes it to the file.
func saveDataURI(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil { // Encode as PNG for data URI
		return "", err
	}
	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())
	dataURI := "data:image/png;base64," + base64Data
	// _, err := file.WriteString(dataURI)
	// return err
	return dataURI, nil
}

func encodeWep(w *os.File, img image.Image) error {
	bounds := img.Bounds()

	// Convert image to RGBA
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Encode with libwebp quality settings
	// Quality: 0-100 (higher is better quality)
	// You can adjust the quality parameter as needed
	quality := float32(90)

	// Create WebP encoder options
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return err
	}

	// Encode to WebP
	return webp.Encode(w, rgba, options)
}
