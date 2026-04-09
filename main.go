package main

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

const mode = ModeSextant

func main() {
	var err = termbox.Init()
	if err != nil {
		fmt.Print("Failed to initialize!")
		return
	}
	termbox.SetInputMode(termbox.InputMouse)
	termbox.SetOutputMode(termbox.Output256)

	var atlas, atlasW, _ = LoadPNG("atlas.png")
	InitTileGrid(mode)

	var eventQueue = make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	var ticks, lastTPS = 0, 0
	var lastTPSTime = time.Now()
	var lastRedrawTime = time.Now()

	var redraw = func() {
		var now = time.Now()
		var elapsed = now.Sub(lastRedrawTime).Seconds()
		lastRedrawTime = now

		var fps int
		if elapsed > 0 {
			fps = int(1.0/elapsed + 0.5)
		}

		var fpsText = appendFPS(fpsBuf[:0], fps)
		fpsText = append(fpsText, "  "...)
		fpsText = appendTPS(fpsText, lastTPS)

		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		SetTile(0, 0, Tile{0, 8, 17}, mode, atlas, atlasW)
		SetTile(1, 0, Tile{1, 2, 52}, mode, atlas, atlasW)
		SetTile(0, 1, Tile{2, 3, 22}, mode, atlas, atlasW)
		SetTile(1, 1, Tile{3, 5, 53}, mode, atlas, atlasW)
		drawString(0, 0, termbox.ColorWhite, termbox.ColorBlack, fpsText)
		drawString(0, 1, termbox.ColorWhite, termbox.ColorBlack, writeMemoryUsage(debugBuf[:0]))
		termbox.Flush()
	}
	redraw()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				termbox.Close()
				return
			} else if ev.Type == termbox.EventResize {
				InitTileGrid(mode)
				redraw()
			}
		default:
		}

		ticks++
		if time.Since(lastTPSTime) >= time.Second {
			lastTPS = ticks
			ticks = 0
			lastTPSTime = time.Now()
			redraw()
		}
	}
}

func drawString(x, y int, fg, bg termbox.Attribute, msg []byte) {
	var startX = x
	for _, c := range msg {
		if c == '\n' {
			y++
			x = startX
			continue
		}
		termbox.SetCell(x, y, rune(c), fg, bg)
		x++
	}
}
