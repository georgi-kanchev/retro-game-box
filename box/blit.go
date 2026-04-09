package box

import "github.com/nsf/termbox-go"

// RenderMode selects the Unicode block character set and sub-cell pixel grid.
type RenderMode int

const (
	ModeSextant RenderMode = iota // 2×3 px/cell -> width=6, height=4
	ModeHalf                      // 1×2 px/cell -> width=12, height=6
	ModeQuarter                   // 2×2 px/cell -> width=6, height=6
	ModeBraille                   // 2×4 px/cell -> width=6, height=3
)

// cellSize returns the terminal cells per tile on each axis for a given mode.
func cellSize(mode RenderMode) (tcw, tch int) {
	switch mode {
	case ModeHalf:
		return 12, 6
	case ModeQuarter:
		return 6, 6
	case ModeBraille:
		return 6, 3
	default: // ModeSextant
		return 6, 4
	}
}

// renderMode dispatches pixels to the correct Blit* function for mode.
func renderMode(mode RenderMode, pixels []termbox.Attribute, width, height, offX, offY int, bg termbox.Attribute) {
	switch mode {
	case ModeHalf:
		BlitHalf(pixels, width, height, offX, offY, bg)
	case ModeQuarter:
		BlitQuarter(pixels, width, height, offX, offY, bg)
	case ModeBraille:
		BlitBraille(pixels, width, height, offX, offY, bg)
	default: // ModeSextant
		Blit(pixels, width, height, offX, offY, bg)
	}
}

// sextantRune converts a 6-bit bitmask to a Unicode block sextant character.
// Bit layout per terminal cell:
//
//	col:  0  1
//	row 0: bit 0, bit 1
//	row 1: bit 2, bit 3
//	row 2: bit 4, bit 5
func sextantRune(bitmask int) rune {
	switch bitmask {
	case 0:
		return ' ' // Empty
	case 21:
		return '▌' // U+258C Left Half Block
	case 42:
		return '▐' // U+2590 Right Half Block
	case 63:
		return '█' // U+2588 Full Block
	default:
		// The Unicode sextant block (U+1FB00–U+1FB3B) skips the above values.
		// We must subtract from the offset to account for the skipped bitmasks.
		var offset = bitmask - 1
		if bitmask > 21 {
			offset-- // Shift down because 21 was skipped
		}
		if bitmask > 42 {
			offset-- // Shift down again because 42 was skipped
		}
		return rune(0x1FB00 + offset)
	}
}

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

// quarterblockRune converts a 4-bit bitmask to a Unicode quadrant block character.
// Bit layout per terminal cell:
//
//	col:  0  1
//	row 0: bit 0, bit 1
//	row 1: bit 2, bit 3
func quarterblockRune(bitmask int) rune {
	switch bitmask {
	case 0:
		return ' '
	case 1:
		return '▘' // U+2598 Upper Left
	case 2:
		return '▝' // U+259D Upper Right
	case 3:
		return '▀' // U+2580 Upper Half
	case 4:
		return '▖' // U+2596 Lower Left
	case 5:
		return '▌' // U+258C Left Half
	case 6:
		return '▞' // U+259E Diagonal (UR+LL)
	case 7:
		return '▛' // U+259B Three-quarter (not LR)
	case 8:
		return '▗' // U+2597 Lower Right
	case 9:
		return '▚' // U+259A Diagonal (UL+LR)
	case 10:
		return '▐' // U+2590 Right Half
	case 11:
		return '▜' // U+259C Three-quarter (not LL)
	case 12:
		return '▄' // U+2584 Lower Half
	case 13:
		return '▙' // U+2599 Three-quarter (not UR)
	case 14:
		return '▟' // U+259F Three-quarter (not UL)
	default:
		return '█' // U+2588 Full
	}
}

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

// BlitQuarter renders a flat row-major pixel array to the termbox backbuffer.
// Each 2×2 block of pixels maps to one terminal cell.
func BlitQuarter(pixels []termbox.Attribute, width, height, offX, offY int, bg termbox.Attribute) {
	var tw, th = termbox.Size()
	var cellCols = (width + 1) / 2
	var cellRows = (height + 1) / 2

	for cy := range cellRows {
		for cx := range cellCols {
			if offX+cx >= tw || offY+cy >= th {
				continue
			}
			var px = cx * 2
			var py = cy * 2
			var bitmask int
			var fg = termbox.ColorDefault
			for row := range 2 {
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
			termbox.SetCell(offX+cx, offY+cy, quarterblockRune(bitmask), fg, bg)
		}
	}
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

// Blit renders a flat row-major pixel array to the termbox backbuffer.
// Each 2×3 block of pixels maps to one terminal cell.
func Blit(pixels []termbox.Attribute, width, height, offX, offY int, bg termbox.Attribute) {
	var tw, th = termbox.Size()
	var cellCols = (width + 1) / 2
	var cellRows = (height + 2) / 3

	for cy := range cellRows {
		for cx := range cellCols {
			if offX+cx >= tw || offY+cy >= th {
				continue
			}
			var px = cx * 2
			var py = cy * 3
			var bitmask int
			var fg = termbox.ColorDefault
			for row := range 3 {
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
			termbox.SetCell(offX+cx, offY+cy, sextantRune(bitmask), fg, bg)
		}
	}
}
