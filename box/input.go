package box

import "github.com/nsf/termbox-go"

// Key represents a keyboard key.
type Key = termbox.Key

// MouseBtn represents a mouse button.
type MouseBtn int

const (
	MouseLeft MouseBtn = iota
	MouseRight
	MouseMiddle
	MouseWheelUp
	MouseWheelDown
)

// Keyboard key constants.
const (
	KeyEsc        = termbox.KeyEsc
	KeyEnter      = termbox.KeyEnter
	KeySpace      = termbox.KeySpace
	KeyBackspace  = termbox.KeyBackspace2
	KeyTab        = termbox.KeyTab
	KeyArrowUp    = termbox.KeyArrowUp
	KeyArrowDown  = termbox.KeyArrowDown
	KeyArrowLeft  = termbox.KeyArrowLeft
	KeyArrowRight = termbox.KeyArrowRight
	KeyDelete     = termbox.KeyDelete
	KeyHome       = termbox.KeyHome
	KeyEnd        = termbox.KeyEnd
	KeyPgUp       = termbox.KeyPgup
	KeyPgDn       = termbox.KeyPgdn
	KeyF1         = termbox.KeyF1
	KeyF2         = termbox.KeyF2
	KeyF3         = termbox.KeyF3
	KeyF4         = termbox.KeyF4
	KeyF5         = termbox.KeyF5
	KeyF6         = termbox.KeyF6
	KeyF7         = termbox.KeyF7
	KeyF8         = termbox.KeyF8
	KeyF9         = termbox.KeyF9
	KeyF10        = termbox.KeyF10
	KeyF11        = termbox.KeyF11
	KeyF12        = termbox.KeyF12
)

var pressedKeys [32]Key
var pressedKeyCount int
var pressedRunes [64]rune
var pressedRuneCount int
var curMouseX, curMouseY int
var pressedMouseBtns [5]bool

func resetInput() {
	pressedKeyCount = 0
	pressedRuneCount = 0
	for i := range pressedMouseBtns {
		pressedMouseBtns[i] = false
	}
}

// processEvent records input from ev into the per-tick input state.
func processEvent(ev termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		if ev.Key != 0 {
			if pressedKeyCount < len(pressedKeys) {
				pressedKeys[pressedKeyCount] = ev.Key
				pressedKeyCount++
			}
		} else if ev.Ch != 0 {
			if pressedRuneCount < len(pressedRunes) {
				pressedRunes[pressedRuneCount] = ev.Ch
				pressedRuneCount++
			}
		}
	case termbox.EventMouse:
		curMouseX = ev.MouseX
		curMouseY = ev.MouseY
		switch ev.Key {
		case termbox.MouseLeft:
			pressedMouseBtns[MouseLeft] = true
		case termbox.MouseRight:
			pressedMouseBtns[MouseRight] = true
		case termbox.MouseMiddle:
			pressedMouseBtns[MouseMiddle] = true
		case termbox.MouseWheelUp:
			pressedMouseBtns[MouseWheelUp] = true
		case termbox.MouseWheelDown:
			pressedMouseBtns[MouseWheelDown] = true
		}
	}
}

// KeyPressed returns true if k was pressed this tick.
func KeyPressed(k Key) bool {
	for i := range pressedKeyCount {
		if pressedKeys[i] == k {
			return true
		}
	}
	return false
}

// RunePressed returns true if r was typed this tick.
func RunePressed(r rune) bool {
	for i := range pressedRuneCount {
		if pressedRunes[i] == r {
			return true
		}
	}
	return false
}

// MouseX returns the mouse column in terminal cells.
func MouseX() int { return curMouseX }

// MouseY returns the mouse row in terminal cells.
func MouseY() int { return curMouseY }

// MousePressed returns true if btn was pressed this tick.
func MousePressed(btn MouseBtn) bool {
	return pressedMouseBtns[btn]
}
