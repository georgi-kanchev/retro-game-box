package box

import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)

// CurrentTPS is the measured ticks per second from the last completed second.
var CurrentTPS int

// CurrentFPS is the measured frames drawn per second from the last completed second.
var CurrentFPS int

var engineMode RenderMode
var dirty bool
var quit bool
var interval time.Duration
var ticks, frames int
var lastSecond time.Time
var eventQueue chan termbox.Event

// Dirty marks the current tick as needing a screen flush.
func Dirty() {
	dirty = true
}

// Quit signals the engine to stop after the current tick.
func Quit() {
	quit = true
}

// Run starts the engine and blocks until completion.
func Run(mode RenderMode, tps int, atlasPath string, update func()) {
	if err := termbox.Init(); err != nil {
		fmt.Println(err)
		return
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputMouse)
	termbox.SetOutputMode(termbox.Output256)

	engineMode = mode
	engineAtlas, engineAtlasW, _ = LoadPNG(atlasPath)
	InitTileGrid(mode)

	interval = time.Second / time.Duration(tps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	eventQueue = make(chan termbox.Event, 16)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	lastSecond = time.Now()

	for range ticker.C {
		now := time.Now()
		ticks++

		resetInput()
		drainEvents()

		// Logic update
		update()

		if dirty {
			termbox.Flush()
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			dirty = false
			frames++
		}

		if now.Sub(lastSecond) >= time.Second {
			CurrentTPS = ticks
			CurrentFPS = frames
			ticks, frames = 0, 0
			lastSecond = now
		}

		if quit {
			return
		}
	}
}

func drainEvents() {
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventResize {
				InitTileGrid(engineMode)
				dirty = true
			}
			processEvent(ev)
		default:
			return
		}
	}
}

func DrawString(x, y int, fg, bg termbox.Attribute, msg []byte) {
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
