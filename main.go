package main

import (
	"retro-game-box/box"
	"retro-game-box/box/input"

	"github.com/gdamore/tcell/v2"
)

var sound *box.Audio

func main() {
	sound, _ = box.LoadSound("music.mp3")

	box.Run(16, 6000, "atlas2.png", update)
}

func update() {
	box.HideCursor()

	for y := range 9 {
		for x := range 16 {
			box.SetTile(x, y, box.Tile{ID: 23, FG: 8, BG: 18})
		}
	}

	if box.KeyIsJustPressedAndHeld(input.LowercaseA) {
		box.PlaySound(sound)
	}

	box.SetTile(1, 1, box.Tile{ID: 25, FG: 2, BG: 22})
	box.SetTile(0, 2, box.Tile{ID: 26, FG: 3, BG: 23})
	box.SetTile(1, 2, box.Tile{ID: 27, FG: 4, BG: 24})

	var statsStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	box.DrawString(0, 0, statsStyle, box.WriteFPS())
	box.DrawString(0, 1, statsStyle, box.WriteMemoryUsage())
}
