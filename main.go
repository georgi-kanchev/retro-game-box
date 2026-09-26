package main

import (
	"retro-game-box/box"
	"time"

	"github.com/nsf/termbox-go"
)

var lastRedraw = time.Now()

func main() {
	box.Run(16, 16, 60, "atlas2.png", update)
}

func update() {
	if box.KeyPressed(box.KeyEsc) {
		box.Quit()
	}

	if time.Since(lastRedraw) >= time.Second {
		lastRedraw = time.Now()
		box.Dirty()
	}

	for y := range 9 {
		for x := range 16 {
			box.SetTile(x, y, box.Tile{ID: 23, FG: 8, BG: 17})
		}
	}
	box.SetTile(1, 1, box.Tile{ID: 25, FG: 2, BG: 20})
	box.SetTile(0, 2, box.Tile{ID: 26, FG: 3, BG: 20})
	box.SetTile(1, 2, box.Tile{ID: 27, FG: 4, BG: 20})

	box.DrawString(0, 0, termbox.ColorWhite, termbox.ColorBlack, box.WriteStats())
	//box.DrawString(0, 1, termbox.ColorWhite, termbox.ColorBlack, box.WriteMemoryUsage())
}
