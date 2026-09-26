package box

import "github.com/gdamore/tcell/v2"

// Key represents a keyboard key.
type Key = tcell.Key

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
	KeyEsc        = tcell.KeyEsc
	KeyEnter      = tcell.KeyEnter
	KeySpace      = tcell.Key(' ')
	KeyBackspace  = tcell.KeyBackspace2
	KeyTab        = tcell.KeyTab
	KeyArrowUp    = tcell.KeyUp
	KeyArrowDown  = tcell.KeyDown
	KeyArrowLeft  = tcell.KeyLeft
	KeyArrowRight = tcell.KeyRight
	KeyDelete     = tcell.KeyDelete
	KeyHome       = tcell.KeyHome
	KeyEnd        = tcell.KeyEnd
	KeyPgUp       = tcell.KeyPgUp
	KeyPgDn       = tcell.KeyPgDn
	KeyF1         = tcell.KeyF1
	KeyF2         = tcell.KeyF2
	KeyF3         = tcell.KeyF3
	KeyF4         = tcell.KeyF4
	KeyF5         = tcell.KeyF5
	KeyF6         = tcell.KeyF6
	KeyF7         = tcell.KeyF7
	KeyF8         = tcell.KeyF8
	KeyF9         = tcell.KeyF9
	KeyF10        = tcell.KeyF10
	KeyF11        = tcell.KeyF11
	KeyF12        = tcell.KeyF12
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
func processEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		if ev.Key() != tcell.KeyRune {
			if pressedKeyCount < len(pressedKeys) {
				pressedKeys[pressedKeyCount] = ev.Key()
				pressedKeyCount++
			}
		} else if ev.Rune() != 0 {
			if pressedRuneCount < len(pressedRunes) {
				pressedRunes[pressedRuneCount] = ev.Rune()
				pressedRuneCount++
			}
		}
	case *tcell.EventMouse:
		curMouseX, curMouseY = ev.Position()
		var btns = ev.Buttons()
		if btns&tcell.Button1 != 0 {
			pressedMouseBtns[MouseLeft] = true
		}
		if btns&tcell.Button2 != 0 {
			pressedMouseBtns[MouseRight] = true
		}
		if btns&tcell.Button3 != 0 {
			pressedMouseBtns[MouseMiddle] = true
		}
		if btns&tcell.WheelUp != 0 {
			pressedMouseBtns[MouseWheelUp] = true
		}
		if btns&tcell.WheelDown != 0 {
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
