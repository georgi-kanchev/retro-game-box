package box

import (
	"time"

	"github.com/nsf/termbox-go"
)

// CurrentTPS is the measured ticks per second from the last completed second.
var CurrentTPS int

// CurrentFPS is the measured frames drawn per second from the last completed second.
var CurrentFPS int

// CurrentIdleTPS is the number of tick-budgets worth of idle time per second.
// It reflects how much headroom is available for computation.
var CurrentIdleTPS int

var engineMode RenderMode
var dirty bool
var quit bool
var tickStart time.Time
var interval time.Duration
var ticks, frames, idleTicks int
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

func idling() bool {
	if time.Since(tickStart) < interval {
		idleTicks++
		return true
	}
	return false
}

// Init initializes termbox, loads the atlas, and sets up the engine.
// Call once before the game loop.
func Init(mode RenderMode, tps int, atlasPath string) error {
	var err = termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetInputMode(termbox.InputMouse)
	termbox.SetOutputMode(termbox.Output256)

	engineMode = mode
	engineAtlas, engineAtlasW, _ = LoadPNG(atlasPath)
	InitTileGrid(mode)

	interval = time.Second / time.Duration(tps)
	eventQueue = make(chan termbox.Event, 16)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	lastSecond = time.Now()
	return nil
}

// Running advances the engine by one tick and returns true while the game
// should keep running. Call it as the condition of a for loop.
//
// On each call it:
//  1. If Dirty() was called: flushes the frame to the terminal, then clears
//     the buffer. Sleeps the remaining tick budget (skipped on first call).
//  2. Drains pending events and updates input state.
//
// Returns false (and closes termbox) when Quit() has been called.
func Running() bool {
	if !tickStart.IsZero() {
		if dirty {
			termbox.Flush()
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			dirty = false
			frames++
		}
		var elapsed = time.Since(tickStart)
		if elapsed < interval {
			time.Sleep(interval - elapsed)
		}
		ticks++
		if time.Since(lastSecond) >= time.Second {
			CurrentTPS = ticks
			CurrentFPS = frames
			CurrentIdleTPS = idleTicks
			ticks, frames, idleTicks = 0, 0, 0
			lastSecond = time.Now()
		}
	}

	tickStart = time.Now()
	resetInput()
	drainEvents()

	if quit {
		termbox.Close()
		return false
	}

	return true
}

func drainEvents() {
drain:
	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventResize {
				InitTileGrid(engineMode)
				dirty = true
			}
			processEvent(ev)
		default:
			break drain
		}
	}
}

// DrawString writes msg to the termbox cell buffer starting at (x, y).
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
