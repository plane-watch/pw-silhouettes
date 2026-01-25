package main

import (
	"build_spritesheet/lib/airframe"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

//go:embed original_sprites.png
var originalSpriteData []byte

const (
	spriteWidth  = 72
	spriteHeight = 72
)

func airframesFromDir(dir string) ([]*airframe.Airframe, error) {
	listing, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %w", err)
	}

	out := make([]*airframe.Airframe, 0, len(listing))

	for _, entry := range listing {
		if entry.IsDir() {
			log.Debug().Str("dir", entry.Name()).Msg("skipping dir")
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".json") {
			log.Debug().Str("file", entry.Name()).Msg("skipping non-json file")
			continue
		}

		af, err := airframe.FromFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to process file: %w", err)
		}

		out = append(out, af)
		log.Info().
			Str("icao", af.ICAO.Designator).
			Msg("added airframe")
	}

	return out, nil
}

func buildSpriteMap(airframes []*airframe.Airframe, idOffset int) map[string]int {
	spriteSet := make(map[string]int)
	n := 0 + idOffset
	for _, af := range airframes {
		for _, frame := range af.Art.Frames {
			if _, ok := spriteSet[frame.Src]; ok {
				log.Info().
					Str("airframe", af.ICAO.Designator).
					Str("src", frame.Src).
					Msg("skipping duplicate")
				continue
			}
			log.Info().
				Str("airframe", af.ICAO.Designator).
				Str("src", frame.Src).
				Int("sprite_id", n).
				Msg("adding sprite")
			spriteSet[frame.Src] = n
			n++
		}
	}
	return spriteSet
}

func runApp(ctx context.Context, cmd *cli.Command) error {

	// read airframe data from json files
	airframes, err := airframesFromDir(cmd.String("airframes_path"))
	if err != nil {
		return err
	}

	// open existing spritesheet
	img, err := png.Decode(bytes.NewBuffer(originalSpriteData))
	if err != nil {
		return fmt.Errorf("failed to decode fallback spritesheet: %w", err)
	}
	bounds := img.Bounds()

	// We will use the original spritesheet width.
	// We will increase the height based on number of new sprites.
	width, height := bounds.Max.X, bounds.Max.Y
	spritesPerRow := width / spriteWidth
	//fmt.Println("sprites per row:", spritesPerRow)
	rows := height / spriteHeight
	//fmt.Println("rows:", rows)
	existingMaxSpriteID := (spritesPerRow * rows) - 1 // -1 as zero indexed

	// generate unique set of sprites (as some airframes reference the same sprites)
	newSprites := buildSpriteMap(airframes, existingMaxSpriteID+1)

	// Work out how much additional height we should add
	numNewSprites := len(newSprites)
	numNewRows := int(math.Ceil(float64(numNewSprites) / float64(spritesPerRow)))
	//fmt.Printf("%d new sprites needs %d new rows added\n", numNewSprites, numNewRows)

	// Work out new img height
	extraHeight := numNewRows * spriteHeight
	//fmt.Println("extra height needed:", extraHeight)
	newHeight := height + extraHeight
	//fmt.Println("new height:", newHeight)

	// Create new image
	newImg := image.NewNRGBA(image.Rect(0, 0, width, newHeight))

	// Copy existing spritesheet into new image
	drawImageOnto(img, newImg, 0, 0)

	for svgFile, spriteNum := range newSprites {
		log.Info().
			Int("sprite_id", spriteNum).
			Str("svg_file", svgFile).
			Msg("adding sprite to spritesheet")

		offX, offY, err := TopLeft(spriteNum, width, spriteWidth, spriteHeight, 0, 0)
		if err != nil {
			return fmt.Errorf("failed to get top left: %w", err)
		}
		err = drawSVGOnto(svgFile, newImg, offX+1, offY+1, cmd.String("inkscape_binary"))
		if err != nil {
			return fmt.Errorf("failed to draw SVG onto new spritesheet: %w", err)
		}
	}

	// Finally, write the new spritesheet
	buf := new(bytes.Buffer)
	err = png.Encode(buf, newImg)
	if err != nil {
		return fmt.Errorf("failed to encode new spritesheet: %w", err)
	}
	err = os.WriteFile(cmd.String("output_png"), buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write new spritesheet: %w", err)
	}

	return nil
}

func drawImageOnto(src, dst image.Image, offsetX, offsetY int) {
	for y := 0; y < src.(*image.NRGBA).Bounds().Dy(); y++ {
		for x := 0; x < src.(*image.NRGBA).Bounds().Dx(); x++ {
			dst.(*image.NRGBA).Set(x+offsetX, y+offsetY, src.At(x, y))
		}
	}
}

func drawSVGOnto(src string, dst image.Image, offsetX, offsetY int, inkscapeBinary string) error {
	tmpDir, err := os.MkdirTemp("", "svg2png-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpPngPath := filepath.Join(tmpDir, "out.png")

	if err := convertSVGtoPNG(inkscapeBinary, src, tmpPngPath); err != nil {
		return fmt.Errorf("failed to convert SVG to PNG: %w", err)
	}

	f, err := os.Open(tmpPngPath)
	if err != nil {
		return fmt.Errorf("failed to open exported png: %w", err)
	}
	defer f.Close()

	pngImage, err := png.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	drawImageOnto(pngImage, dst, offsetX, offsetY)
	return nil
}

// TopLeft returns the (x,y) pixel coords of the top-left corner of the sprite
// in a uniform grid spritesheet.
//
// index: 0-based sprite index (left-to-right, then top-to-bottom)
// sheetW: spritesheet width in pixels
// frameW/frameH: frame size in pixels
// margin: pixels before the first frame (both x and y)
// padding: pixels between frames (both x and y)
func TopLeft(index, sheetW, frameW, frameH, margin, padding int) (x, y int, err error) {
	if index < 0 {
		return 0, 0, fmt.Errorf("index must be >= 0")
	}
	if sheetW <= 0 || frameW <= 0 || frameH <= 0 {
		return 0, 0, fmt.Errorf("sheetW/frameW/frameH must be > 0")
	}
	usableW := sheetW - 2*margin
	if usableW <= 0 {
		return 0, 0, fmt.Errorf("sheetW too small for margin")
	}

	// How many columns fit, accounting for padding between frames.
	cols := (usableW + padding) / (frameW + padding)
	if cols <= 0 {
		return 0, 0, fmt.Errorf("no columns fit (check sheetW/frameW/margin/padding)")
	}

	col := index % cols
	row := index / cols

	x = margin + col*(frameW+padding)
	y = margin + row*(frameH+padding)
	return x, y, nil
}
