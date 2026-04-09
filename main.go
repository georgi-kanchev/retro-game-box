package main

import (
	"fmt"
	"time"

	"retro-game-box/box"

	"github.com/nsf/termbox-go"
)

const targetTPS = 60

var fpsBuf [64]byte
var lastRedraw = time.Now()

func main() {
	if err := box.Init(box.ModeSextant, targetTPS, "atlas.png"); err != nil {
		fmt.Println(err)
		return
	}

	for box.Running() {
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

		var b = fpsBuf[:0]
		b = box.AppendFPS(b, box.CurrentFPS)
		b = append(b, "  "...)
		b = box.AppendTPS(b, box.CurrentTPS)
		b = append(b, "  "...)
		b = box.AppendIdleTPS(b)
		box.DrawString(0, 0, termbox.ColorWhite, termbox.ColorBlack, b)
		box.DrawString(0, 1, termbox.ColorWhite, termbox.ColorBlack, box.WriteMemoryUsage())
	}
}
