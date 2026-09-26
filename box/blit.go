package box

import "github.com/nsf/termbox-go"

// halfblockRune converts a 2-bit bitmask to a Unicode half-block character.
// Bit layout per terminal cell:
//
//	row 0: bit 0  (top)
//	row 1: bit 1  (bottom)
func halfblockRune(bitmask int) rune {
	switch bitmask {
	case 0:
		return ' ' // Empty
	case 1:
		return '▀' // U+2580 Upper Half Block
	case 2:
		return '▄' // U+2584 Lower Half Block
	default: // 3
		return '█' // U+2588 Full Block
	}
}

// BlitHalf renders a flat row-major pixel array to the termbox backbuffer.
// Each 1×2 block of pixels maps to one terminal cell.
func BlitHalf(pixels []termbox.Attribute, width, height, offX, offY int, bg termbox.Attribute) {
	var tw, th = termbox.Size()
	var cellRows = (height + 1) / 2

	for cy := range cellRows {
		for cx := range width {
			if offX+cx >= tw || offY+cy >= th {
				continue
			}
			var py = cy * 2
			var bitmask int
			var fg = termbox.ColorDefault
			for row := range 2 {
				var y = py + row
				if y >= height {
					continue
				}
				var c = pixels[y*width+cx]
				if c != termbox.ColorDefault {
					bitmask |= 1 << row
					if fg == termbox.ColorDefault {
						fg = c
					}
				}
			}
			termbox.SetCell(offX+cx, offY+cy, halfblockRune(bitmask), fg, bg)
		}
	}
}
