package box

import "github.com/gdamore/tcell/v2"

// octantTable maps a 2×4 row-major 8-bit mask (0x00..0xFF) to its precise Unicode character.
//
// Bit layout per terminal cell:
//
//	row 0: bit 0 (left), bit 1 (right)
//	row 1: bit 2 (left), bit 3 (right)
//	row 2: bit 4 (left), bit 5 (right)
//	row 3: bit 6 (left), bit 7 (right)
//
// 230 of the 256 patterns are encoded in the Unicode 16.0 Block Octants range
// U+1CD00..U+1CDE5. The remaining 26 patterns are unified by Unicode with
// pre-existing block characters (Block Elements U+2580..U+259F, Symbols for
// Legacy Computing U+1FB82/U+1FB85, and the quarter blocks U+1CEA0..U+1CEAB,
// U+1FBE6..U+1FBE7); the empty pattern maps to space.
var octantTable = [256]rune{
	/* 0x00 */ ' ', '\U0001cea8', '\U0001ceab', '\U0001fb82', '\U0001cd00', '\U00002598', '\U0001cd01', '\U0001cd02',
	'\U0001cd03', '\U0001cd04', '\U0000259d', '\U0001cd05', '\U0001cd06', '\U0001cd07', '\U0001cd08', '\U00002580',
	/* 0x10 */ '\U0001cd09', '\U0001cd0a', '\U0001cd0b', '\U0001cd0c', '\U0001fbe6', '\U0001cd0d', '\U0001cd0e', '\U0001cd0f',
	'\U0001cd10', '\U0001cd11', '\U0001cd12', '\U0001cd13', '\U0001cd14', '\U0001cd15', '\U0001cd16', '\U0001cd17',
	/* 0x20 */ '\U0001cd18', '\U0001cd19', '\U0001cd1a', '\U0001cd1b', '\U0001cd1c', '\U0001cd1d', '\U0001cd1e', '\U0001cd1f',
	'\U0001fbe7', '\U0001cd20', '\U0001cd21', '\U0001cd22', '\U0001cd23', '\U0001cd24', '\U0001cd25', '\U0001cd26',
	/* 0x30 */ '\U0001cd27', '\U0001cd28', '\U0001cd29', '\U0001cd2a', '\U0001cd2b', '\U0001cd2c', '\U0001cd2d', '\U0001cd2e',
	'\U0001cd2f', '\U0001cd30', '\U0001cd31', '\U0001cd32', '\U0001cd33', '\U0001cd34', '\U0001cd35', '\U0001fb85',
	/* 0x40 */ '\U0001cea3', '\U0001cd36', '\U0001cd37', '\U0001cd38', '\U0001cd39', '\U0001cd3a', '\U0001cd3b', '\U0001cd3c',
	'\U0001cd3d', '\U0001cd3e', '\U0001cd3f', '\U0001cd40', '\U0001cd41', '\U0001cd42', '\U0001cd43', '\U0001cd44',
	/* 0x50 */ '\U00002596', '\U0001cd45', '\U0001cd46', '\U0001cd47', '\U0001cd48', '\U0000258c', '\U0001cd49', '\U0001cd4a',
	'\U0001cd4b', '\U0001cd4c', '\U0000259e', '\U0001cd4d', '\U0001cd4e', '\U0001cd4f', '\U0001cd50', '\U0000259b',
	/* 0x60 */ '\U0001cd51', '\U0001cd52', '\U0001cd53', '\U0001cd54', '\U0001cd55', '\U0001cd56', '\U0001cd57', '\U0001cd58',
	'\U0001cd59', '\U0001cd5a', '\U0001cd5b', '\U0001cd5c', '\U0001cd5d', '\U0001cd5e', '\U0001cd5f', '\U0001cd60',
	/* 0x70 */ '\U0001cd61', '\U0001cd62', '\U0001cd63', '\U0001cd64', '\U0001cd65', '\U0001cd66', '\U0001cd67', '\U0001cd68',
	'\U0001cd69', '\U0001cd6a', '\U0001cd6b', '\U0001cd6c', '\U0001cd6d', '\U0001cd6e', '\U0001cd6f', '\U0001cd70',
	/* 0x80 */ '\U0001cea0', '\U0001cd71', '\U0001cd72', '\U0001cd73', '\U0001cd74', '\U0001cd75', '\U0001cd76', '\U0001cd77',
	'\U0001cd78', '\U0001cd79', '\U0001cd7a', '\U0001cd7b', '\U0001cd7c', '\U0001cd7d', '\U0001cd7e', '\U0001cd7f',
	/* 0x90 */ '\U0001cd80', '\U0001cd81', '\U0001cd82', '\U0001cd83', '\U0001cd84', '\U0001cd85', '\U0001cd86', '\U0001cd87',
	'\U0001cd88', '\U0001cd89', '\U0001cd8a', '\U0001cd8b', '\U0001cd8c', '\U0001cd8d', '\U0001cd8e', '\U0001cd8f',
	/* 0xA0 */ '\U00002597', '\U0001cd90', '\U0001cd91', '\U0001cd92', '\U0001cd93', '\U0000259a', '\U0001cd94', '\U0001cd95',
	'\U0001cd96', '\U0001cd97', '\U00002590', '\U0001cd98', '\U0001cd99', '\U0001cd9a', '\U0001cd9b', '\U0000259c',
	/* 0xB0 */ '\U0001cd9c', '\U0001cd9d', '\U0001cd9e', '\U0001cd9f', '\U0001cda0', '\U0001cda1', '\U0001cda2', '\U0001cda3',
	'\U0001cda4', '\U0001cda5', '\U0001cda6', '\U0001cda7', '\U0001cda8', '\U0001cda9', '\U0001cdaa', '\U0001cdab',
	/* 0xC0 */ '\U00002582', '\U0001cdac', '\U0001cdad', '\U0001cdae', '\U0001cdaf', '\U0001cdb0', '\U0001cdb1', '\U0001cdb2',
	'\U0001cdb3', '\U0001cdb4', '\U0001cdb5', '\U0001cdb6', '\U0001cdb7', '\U0001cdb8', '\U0001cdb9', '\U0001cdba',
	/* 0xD0 */ '\U0001cdbb', '\U0001cdbc', '\U0001cdbd', '\U0001cdbe', '\U0001cdbf', '\U0001cdc0', '\U0001cdc1', '\U0001cdc2',
	'\U0001cdc3', '\U0001cdc4', '\U0001cdc5', '\U0001cdc6', '\U0001cdc7', '\U0001cdc8', '\U0001cdc9', '\U0001cdca',
	/* 0xE0 */ '\U0001cdcb', '\U0001cdcc', '\U0001cdcd', '\U0001cdce', '\U0001cdcf', '\U0001cdd0', '\U0001cdd1', '\U0001cdd2',
	'\U0001cdd3', '\U0001cdd4', '\U0001cdd5', '\U0001cdd6', '\U0001cdd7', '\U0001cdd8', '\U0001cdd9', '\U0001cdda',
	/* 0xF0 */ '\U00002584', '\U0001cddb', '\U0001cddc', '\U0001cddd', '\U0001cdde', '\U00002599', '\U0001cddf', '\U0001cde0',
	'\U0001cde1', '\U0001cde2', '\U0000259f', '\U0001cde3', '\U00002586', '\U0001cde4', '\U0001cde5', '\U00002588',
}

// BlitOctants renders a flat row-major pixel array to the tcell backbuffer.
// Each 2×4 block of pixels maps to one terminal cell as a block octant pattern.
func BlitOctants(pixels []tcell.Color, width, height, offX, offY int, bg tcell.Color) {
	var tw, th = screen.Size()
	var cellCols, cellRows = (width + 1) / 2, (height + 3) / 4

	for cy := range cellRows {
		for cx := range cellCols {
			if offX+cx >= tw || offY+cy >= th {
				continue
			}
			var px, py = cx * 2, cy * 4
			var bitmask uint8
			var fg = tcell.ColorDefault

			for row := range 4 {
				for col := range 2 {
					var x, y = px + col, py + row
					if x >= width || y >= height {
						continue
					}
					var c = pixels[y*width+x]
					if c != tcell.ColorDefault {
						bitmask |= 1 << (row*2 + col)
						if fg == tcell.ColorDefault {
							fg = c
						}
					}
				}
			}
			screen.SetContent(offX+cx, offY+cy, octantTable[bitmask], nil, tcell.StyleDefault.Foreground(fg).Background(bg))
		}
	}
}
