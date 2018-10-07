// barnsley.go creates a gif of a barnsley fern
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math/rand"
	"os"
)

func init() {
	// initialize a random seed
	rand.Seed(42)
}

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0
	blackIndex = 1
)

func main() {
	drawBarnsley(os.Stdout)
}

func drawBarnsley(out io.Writer) {
	const (
		size     = 2200    // image canvas covers [-size..+size]
		scaleX   = 600     // scaling factor for x coordinate
		scaleY   = 600     // scaling factor for y coordinate
		nframes  = 1       // number of animation frames
		delay    = 35      // delay between frames
		treeSize = 5000000 // eventual tree size
	)
	anim := gif.GIF{LoopCount: nframes}
	x := 0.0 // initial x coordinate
	y := 0.0 // initial y coordinate
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for j := 0; j < treeSize; j++ {
			r := rand.Float64()
			x, y = calcBarnsley(r, x, y) // calculate the next x, y coordinates
			err := drawWithStroke(img, int(x*scaleX+(size)), int(y*scaleY), blackIndex, 2)
			if err != nil {
				log.Fatalf("error while drawing points: %v\n", err)
			}
		}
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	if err := gif.EncodeAll(out, &anim); err != nil {
		log.Fatalf("error encoding gif: %v\n", err)
	}
}

func calcBarnsley(number float64, x, y float64) (float64, float64) {
	var a float64
	var b float64
	switch {
	case number <= 0.01:
		a = 0
		b = 0.16 * y
		return a, b
	case number <= 0.86:
		a = 0.85*x + 0.04*y
		b = -0.04*x + 0.85*y + 1.6
		return a, b
	case number <= 0.93:
		a = 0.2*x - 0.26*y
		b = 0.23*x + 0.22*y + 1.6
		return a, b
	case number <= 1.0:
		a = -0.15*x + 0.28*y
		b = 0.26*x + 0.24*y + 0.44
		return a, b
	}
	return a, b
}

// only a value of `2' for `stroke' is supported
func drawWithStroke(img *image.Paletted, x int, y int, idx uint8, stroke int) error {
	switch stroke {
	case 1:
		img.SetColorIndex(x, y, idx)
	case 2:
		img.SetColorIndex(x, y, idx)
		img.SetColorIndex(x+1, y, idx)
		img.SetColorIndex(x, y+1, idx)
		img.SetColorIndex(x+1, y+1, idx)
		img.SetColorIndex(x-1, y, idx)
		img.SetColorIndex(x, y-1, idx)
		img.SetColorIndex(x-1, y-1, idx)
		img.SetColorIndex(x+1, y-1, idx)
		img.SetColorIndex(x-1, y+1, idx)
	default:
		return fmt.Errorf("stroke of %d is not supported", stroke)
	}
	return nil
}
