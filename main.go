package main

import (
	"retro-game-box/box"
	"time"

	"github.com/nsf/termbox-go"
)

var lastRedraw = time.Now()

func main() {
	box.Run(box.ModeHalf, 1000, "atlas.png", update)
}

func update() {
	if box.KeyPressed(box.KeyEsc) {
		box.Quit()
	}

	if time.Since(lastRedraw) >= time.Second {
		lastRedraw = time.Now()
		box.Dirty()
	}

	box.SetTile(0, 0, box.Tile{ID: 0, FG: 8, BG: 17})
	box.SetTile(1, 0, box.Tile{ID: 1, FG: 2, BG: 52})
	box.SetTile(0, 1, box.Tile{ID: 2, FG: 3, BG: 22})
	box.SetTile(1, 1, box.Tile{ID: 3, FG: 5, BG: 53})

	box.DrawString(0, 0, termbox.ColorWhite, termbox.ColorBlack, box.WriteStats())
	box.DrawString(0, 1, termbox.ColorWhite, termbox.ColorBlack, box.WriteMemoryUsage())
}
