package box

import "github.com/gdamore/tcell/v2"

// Tile identifies a sprite and its colors.
// ID is the 1D index of the tile in the atlas (row-major, zero-based).
type Tile struct {
	ID int
	FG byte
	BG byte
}

var engineAtlas []tcell.Color
var engineAtlasW int

var engineTileSize int

// tilePixels is a reusable scratch buffer for one tile's pixel data.
// It grows as needed to fit the current tile size.
var tilePixels []tcell.Color

// tileGrid is a flat row-major backing store for the tile grid.
// tileGridW is its width in tiles; height is len(tileGrid)/tileGridW.
var tileGrid []Tile
var tileGridW int

// paletteColor converts a 256-color palette index to a tcell color.
// Index 0 means the terminal default color.
func paletteColor(index byte) tcell.Color {
	if index == 0 {
		return tcell.ColorDefault
	}
	return tcell.PaletteColor(int(index))
}

// cellsPerTile returns the terminal cells a tile occupies on each axis.
// Block octants pack 2×4 pixels into each cell.
func cellsPerTile() (tcw, tch int) {
	return (engineTileSize + 1) / 2, (engineTileSize + 3) / 4
}

// InitTileGrid sizes the tile grid to match the current terminal dimensions.
// Call on startup and after every resize.
func InitTileGrid() {
	var tw, th = screen.Size()
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
	var tilesPerRow = engineAtlasW / engineTileSize
	if tilesPerRow <= 0 {
		return
	}
	var srcX = (t.ID % tilesPerRow) * engineTileSize
	var srcY = (t.ID / tilesPerRow) * engineTileSize

	var need = engineTileSize * engineTileSize
	if cap(tilePixels) < need {
		tilePixels = make([]tcell.Color, need)
	} else {
		tilePixels = tilePixels[:need]
	}

	for y := range engineTileSize {
		for x := range engineTileSize {
			if engineAtlas[(srcY+y)*engineAtlasW+(srcX+x)] != tcell.ColorDefault {
				tilePixels[y*engineTileSize+x] = paletteColor(t.FG)
			} else {
				tilePixels[y*engineTileSize+x] = tcell.ColorDefault
			}
		}
	}

	BlitOctants(tilePixels, engineTileSize, engineTileSize, col*tcw, row*tch, paletteColor(t.BG))

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
			screen.SetContent(col*tcw+cx, row*tch+cy, ' ', nil, tcell.StyleDefault)
		}
	}
	var idx = row*tileGridW + col
	if idx >= 0 && idx < len(tileGrid) {
		tileGrid[idx] = Tile{}
	}
}
