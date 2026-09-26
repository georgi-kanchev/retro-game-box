package box

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

// CurrentTPS is the measured ticks per second from the last completed second.
var CurrentTPS int

// CurrentFPS is the measured frames drawn per second from the last completed second.
var CurrentFPS int

var screen tcell.Screen
var dirty bool
var quit bool
var needSync bool
var interval time.Duration
var ticks, frames int
var lastSecond time.Time
var eventQueue chan tcell.Event

// Screen returns the active tcell screen, or nil if the engine is not running.
func Screen() tcell.Screen {
	return screen
}

// Dirty marks the current tick as needing a screen flush.
func Dirty() {
	dirty = true
}

// Quit signals the engine to stop after the current tick.
func Quit() {
	quit = true
}

// HideCursor hides the terminal cursor.
func HideCursor() {
	if screen != nil {
		screen.HideCursor()
	}
}

// ShowCursor shows the terminal cursor at the given cell position.
func ShowCursor(x, y int) {
	if screen != nil {
		screen.ShowCursor(x, y)
	}
}

// Run starts the engine and blocks until completion.
// tileW and tileH are the pixel dimensions of each tile in the atlas.
func Run(tileW, tileH, tps int, atlasPath string, update func()) {
	if tileW <= 0 || tileH <= 0 {
		fmt.Println("box: tile size must be positive")
		return
	}
	var s, err = tcell.NewScreen()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = s.Init(); err != nil {
		fmt.Println(err)
		return
	}
	screen = s
	defer func() {
		screen.Fini()
		screen = nil
	}()

	screen.EnableMouse()

	engineTileW = tileW
	engineTileH = tileH
	engineAtlas, engineAtlasW, _ = LoadPNG(atlasPath)
	InitTileGrid()

	interval = time.Second / time.Duration(tps)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	eventQueue = make(chan tcell.Event, 16)
	go func() {
		for {
			eventQueue <- screen.PollEvent()
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
			if needSync {
				screen.Sync()
				needSync = false
			} else {
				screen.Show()
			}
			screen.Clear()
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
			if ev == nil {
				continue
			}
			if _, ok := ev.(*tcell.EventResize); ok {
				InitTileGrid()
				dirty = true
				needSync = true
			}
			processEvent(ev)
		default:
			return
		}
	}
}

func DrawString(x, y int, style tcell.Style, msg []byte) {
	var startX = x
	for _, c := range msg {
		if c == '\n' {
			y++
			x = startX
			continue
		}
		screen.SetContent(x, y, rune(c), nil, style)
		x++
	}
}
