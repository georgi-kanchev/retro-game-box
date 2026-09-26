package box

import "github.com/nsf/termbox-go"

// brailleRune converts an 8-bit row-major bitmask to a Unicode braille pattern character.
// Bit layout per terminal cell:
//
//	col:  0  1
//	row 0: bit 0, bit 1
//	row 1: bit 2, bit 3
//	row 2: bit 4, bit 5
//	row 3: bit 6, bit 7
//
// All 256 braille patterns (U+2800–U+28FF) are in Unicode 1.0.
func brailleRune(bitmask int) rune {
	// Remap row-major bit positions to braille dot bit positions.
	// Braille: dots 1–3 on col 0 rows 0–2, dots 4–6 on col 1 rows 0–2, dots 7–8 on row 3.
	var dotPos = [8]int{0, 3, 1, 4, 2, 5, 6, 7}
	var b int
	for i, dot := range dotPos {
		if bitmask&(1<<i) != 0 {
			b |= 1 << dot
		}
	}
	return rune(0x2800 + b)
}

// BlitBraille renders a flat row-major pixel array to the termbox backbuffer.
// Each 2×4 block of pixels maps to one terminal cell as a braille pattern.
func BlitBraille(pixels []termbox.Attribute, width, height, offX, offY int, bg termbox.Attribute) {
	var tw, th = termbox.Size()
	var cellCols = (width + 1) / 2
	var cellRows = (height + 3) / 4

	for cy := range cellRows {
		for cx := range cellCols {
			if offX+cx >= tw || offY+cy >= th {
				continue
			}
			var px = cx * 2
			var py = cy * 4
			var bitmask int
			var fg = termbox.ColorDefault
			for row := range 4 {
				for col := range 2 {
					var x, y = px + col, py + row
					if x >= width || y >= height {
						continue
					}
					var c = pixels[y*width+x]
					if c != termbox.ColorDefault {
						bitmask |= 1 << (row*2 + col)
						if fg == termbox.ColorDefault {
							fg = c
						}
					}
				}
			}
			termbox.SetCell(offX+cx, offY+cy, brailleRune(bitmask), fg, bg)
		}
	}
}
