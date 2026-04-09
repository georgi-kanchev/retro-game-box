package main

import "github.com/nsf/termbox-go"

// TileSize is the pixel width and height of each tile.
// 12 divides evenly by all supported cell grids: 1, 2, 3, and 4.
const TileSize = 12

// Tile identifies a sprite and its colors.
// ID is the 1D index of the tile in the atlas (row-major, zero-based).
type Tile struct {
	ID int
	FG byte
	BG byte
}

// tilePixels is a reusable scratch buffer for one tile's pixel data.
var tilePixels [TileSize * TileSize]termbox.Attribute

// tileGrid is a flat row-major backing store for the tile grid.
// tileGridW is its width in tiles; height is len(tileGrid)/tileGridW.
var tileGrid []Tile
var tileGridW int

// InitTileGrid sizes the tile grid to match the current terminal dimensions.
// Call on startup and after every resize. Reuses the existing allocation if
// capacity allows; only allocates when the terminal has grown beyond it.
func InitTileGrid(mode RenderMode) {
	var tw, th = termbox.Size()
	var tcw, tch = cellSize(mode)
	var w, h = tw / tcw, th / tch
	var need = w * h
	if need > cap(tileGrid) {
		tileGrid = make([]Tile, need)
	} else {
		tileGrid = tileGrid[:need]
		clear(tileGrid)
	}
	tileGridW = w
}

// SetTile blits tile t from the atlas at tile grid position (col, row).
// Atlas tiles are row-major, packed at TileSize×TileSize with no spacing.
// Lit atlas pixels (non-ColorDefault) are drawn with color; dark pixels are transparent.
func SetTile(col, row int, t Tile, mode RenderMode, atlas []termbox.Attribute, atlasW int) {
	var tcw, tch = cellSize(mode)
	var tilesPerRow = atlasW / TileSize
	var srcX = (t.ID % tilesPerRow) * TileSize
	var srcY = (t.ID / tilesPerRow) * TileSize

	for y := range TileSize {
		for x := range TileSize {
			if atlas[(srcY+y)*atlasW+(srcX+x)] != termbox.ColorDefault {
				tilePixels[y*TileSize+x] = termbox.Attribute(t.FG)
			} else {
				tilePixels[y*TileSize+x] = termbox.ColorDefault
			}
		}
	}

	renderMode(mode, tilePixels[:], TileSize, TileSize, col*tcw, row*tch, termbox.Attribute(t.BG))

	var idx = row*tileGridW + col
	if idx >= 0 && idx < len(tileGrid) {
		tileGrid[idx] = t
	}
}

// GetTile returns the tile at (col, row) from the backing grid.
func GetTile(col, row int) Tile {
	var idx = row*tileGridW + col
	if idx >= 0 && idx < len(tileGrid) {
		return tileGrid[idx]
	}
	return Tile{}
}

// ClearTile erases the tile at (col, row).
func ClearTile(col, row int, mode RenderMode) {
	var tcw, tch = cellSize(mode)
	for cy := range tch {
		for cx := range tcw {
			termbox.SetCell(col*tcw+cx, row*tch+cy, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	var idx = row*tileGridW + col
	if idx >= 0 && idx < len(tileGrid) {
		tileGrid[idx] = Tile{}
	}
}
