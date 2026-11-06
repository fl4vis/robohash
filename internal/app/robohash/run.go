// Package robohash implements the command-line application.
package robohash

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"os"
	"slices"
	"strings"

	"github.com/fl4vis/robohash/internal/assets"
	"github.com/fl4vis/robohash/internal/generator"
	"github.com/fl4vis/robohash/internal/output"
)

const defaultFormat = "png"

// Run parses command-line flags, generates an image, and writes it to disk.
func Run() error {
	flags := flag.NewFlagSet("robohash", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	input := flags.String("input", "", "String to be hashed")

	hashSlots := flags.Int("slots", 11, "Number of slots hash will be divided")

	ignoreExtension := flags.Bool("ignoreExt", true, "Removes any extension of the input before hashing")

	avatarType := flags.String("type", "robo", `Avatar type: [robo, onepiece] 
		(onepiece doesn't support: [set, color])`)

	set := flags.String("set", "any", "RoboHash set: [1,5] or any")

	color := flags.String("color", "", `[blue, brown, green, grey, orange, pink, purple, red, white, yellow]
	  	If omitted or invalid, a random color is set`)

	backgroundSet := flags.String("bgset", "", "Background: [1,2] or any")

	format := flags.String("extension", defaultFormat, "Output format: [png, jpg, jpeg, gif, ppm, datauri]")

	width := flags.Int("sizex", 300, "Output width in pixels")
	height := flags.Int("sizey", 300, "Output height in pixels")

	flags.Usage = func() {
		banner := `
██████╗  ██████╗ ██████╗  ██████╗ ██╗  ██╗ █████╗ ███████╗██╗  ██╗
██╔══██╗██╔═══██╗██╔══██╗██╔═══██╗██║  ██║██╔══██╗██╔════╝██║  ██║
██████╔╝██║   ██║██████╔╝██║   ██║███████║███████║███████╗███████║
██╔══██╗██║   ██║██╔══██╗██║   ██║██╔══██║██╔══██║╚════██║██╔══██║
██║  ██║╚██████╔╝██████╔╝╚██████╔╝██║  ██║██╔══██║███████║██║  ██║
╚═╝  ╚═════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝

Usage:
`
		fmt.Fprint(flags.Output(), banner)
		flags.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(
				flags.Output(),
				"  \033[36m%-12s\033[0m %s\n",
				"-"+f.Name,
				f.Usage,
			)
		})
	}
	if err := flags.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if *input == "" {
		return fmt.Errorf("input value is required")
	}
	*avatarType = strings.ToLower(*avatarType)
	*format = strings.ToLower(*format)
	if *width < 1 || *height < 1 {
		return fmt.Errorf("image dimensions must be positive")
	}
	if !validFormat(*format) {
		return fmt.Errorf("unsupported extension %q", *format)
	}

	img, err := generate(*avatarType, *input, *hashSlots, *ignoreExtension, *set, *color, *backgroundSet, *format, *width, *height)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s.%s", *input, *format)
	return output.SaveImage(filename, *format, img, os.Stdout)
}

func generate(avatarType, input string, slots int, ignoreExtension bool, set, color, background, format string, width, height int) (image.Image, error) {
	switch avatarType {
	case "robo":
		return generator.GenerateRobo(assets.Files, input, slots, ignoreExtension, generator.RoboOptions{
			Set:        set,
			Color:      color,
			Background: background,
			Format:     format,
			Width:      width,
			Height:     height,
		})
	case "onepiece":
		return generator.GenerateOnePiece(assets.Files, input, slots, ignoreExtension, generator.OnePieceOptions{
			Background: background,
			Format:     format,
			Width:      width,
			Height:     height,
		})
	default:
		return nil, fmt.Errorf("unknown avatar type %q; use robo or onepiece", avatarType)
	}
}

func validFormat(format string) bool {
	return slices.Contains([]string{"png", "jpg", "jpeg", "gif", "ppm", "datauri"}, format)
}
