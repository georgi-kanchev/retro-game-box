package box

import "github.com/nsf/termbox-go"

// Tile identifies a sprite and its colors.
// ID is the 1D index of the tile in the atlas (row-major, zero-based).
type Tile struct {
	ID int
	FG byte
	BG byte
}

var engineAtlas []termbox.Attribute
var engineAtlasW int

// engineTileW and engineTileH are the pixel dimensions of each tile.
// They are set by Run and may differ between runs.
var engineTileW, engineTileH int

// tilePixels is a reusable scratch buffer for one tile's pixel data.
// It grows as needed to fit the current tile size.
var tilePixels []termbox.Attribute

// tileGrid is a flat row-major backing store for the tile grid.
// tileGridW is its width in tiles; height is len(tileGrid)/tileGridW.
var tileGrid []Tile
var tileGridW int

// cellsPerTile returns the terminal cells a tile occupies on each axis.
// Block octants pack 2×4 pixels into each cell.
func cellsPerTile() (tcw, tch int) {
	return (engineTileW + 1) / 2, (engineTileH + 3) / 4
}

// InitTileGrid sizes the tile grid to match the current terminal dimensions.
// Call on startup and after every resize.
func InitTileGrid() {
	var tw, th = termbox.Size()
	var tcw, tch = cellsPerTile()
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
func SetTile(col, row int, t Tile) {
	var tcw, tch = cellsPerTile()
	var tilesPerRow = engineAtlasW / engineTileW
	if tilesPerRow <= 0 {
		return
	}
	var srcX = (t.ID % tilesPerRow) * engineTileW
	var srcY = (t.ID / tilesPerRow) * engineTileH

	var need = engineTileW * engineTileH
	if cap(tilePixels) < need {
		tilePixels = make([]termbox.Attribute, need)
	} else {
		tilePixels = tilePixels[:need]
	}

	for y := range engineTileH {
		for x := range engineTileW {
			if engineAtlas[(srcY+y)*engineAtlasW+(srcX+x)] != termbox.ColorDefault {
				tilePixels[y*engineTileW+x] = termbox.Attribute(t.FG)
			} else {
				tilePixels[y*engineTileW+x] = termbox.ColorDefault
			}
		}
	}

	BlitOctants(tilePixels, engineTileW, engineTileH, col*tcw, row*tch, termbox.Attribute(t.BG))

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
func ClearTile(col, row int) {
	var tcw, tch = cellsPerTile()
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
