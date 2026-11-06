package generator

import (
	"fmt"
	"image"
	"io/fs"
)

const onePieceRoot = "img/one-piece"

type OnePieceOptions struct {
	Format     string
	Width      int
	Height     int
	Background string
}

func GenerateOnePiece(fsys fs.FS, input string, slots int, ignoreExtension bool, options OnePieceOptions) (image.Image, error) {
	hash, err := NewHashSlots(input, slots, ignoreExtension)
	if err != nil {
		return nil, err
	}

	path, err := selectedFile(fsys, onePieceRoot, hash.At(0))
	if err != nil {
		return nil, fmt.Errorf("select One Piece image: %w", err)
	}

	backgroundSets, err := directories(fsys, roboRoot+"/backgrounds")
	if err != nil {
		return nil, err
	}

	backgroundSet := selectBackground(options.Background, backgroundSets, hash)

	img, err := decodeImage(fsys, path)
	if err != nil {
		return nil, err
	}
	result := resizeWorkingImage(img)

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
