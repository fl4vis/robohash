// Package output encodes generated images.
package output

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

// SaveImage writes an image in the requested format. Data URIs are written to stdout.
func SaveImage(filename, format string, img image.Image, stdout io.Writer) error {
	if format == "datauri" {
		dataURI, err := dataURI(img)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(stdout, dataURI)
		return err
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
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case "gif":
		return gif.Encode(file, img, nil)
	case "ppm":
		return encodePPM(file, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func encodePPM(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	if _, err := fmt.Fprintf(w, "P6\n%d %d\n255\n", bounds.Dx(), bounds.Dy()); err != nil {
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

func dataURI(img image.Image) (string, error) {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}
