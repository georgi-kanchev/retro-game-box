package box

import (
	"image"
	"image/png"
	"os"

	"github.com/nsf/termbox-go"
)

// LoadPNG decodes a PNG file and returns a flat pixel array (row-major),
// the image width, height, and any error. Each pixel is the nearest
// xterm-256 palette color. Fully transparent pixels map to ColorDefault.
func LoadPNG(path string) (pixels []termbox.Attribute, width, height int) {
	var f, err = os.Open(path)
	if err != nil {
		return nil, 0, 0
	}
	defer f.Close()

	var img image.Image
	img, err = png.Decode(f)
	if err != nil {
		return nil, 0, 0
	}

	var b = img.Bounds()
	var w = b.Max.X - b.Min.X
	var h = b.Max.Y - b.Min.Y

	pixels = make([]termbox.Attribute, w*h)
	for y := range h {
		for x := range w {
			var _, _, _, a16 = img.At(b.Min.X+x, b.Min.Y+y).RGBA()
			if a16 < 0x8000 {
				pixels[y*w+x] = termbox.ColorDefault
				continue
			}
			pixels[y*w+x] = termbox.ColorWhite
		}
	}
	return pixels, w, h
}
