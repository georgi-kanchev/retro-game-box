package box

import (
	"retro-game-box/box/input"

	"github.com/gdamore/tcell/v2"
)

// Whether the key produced an input event this tick. It fires on the initial press and on every OS auto-repeat.
func KeyIsJustPressedAndHeld(k input.Key) bool {
	for _, i := range currentKeys {
		if i == k {
			return true
		}
	}
	return false
}

// Whether the mouse button was pressed this tick. For the wheel, it reports a scroll tick.
func MouseIsJustPressed(b input.Button) bool { return justPressedButtons[b] }

// Whether the mouse button is currently held down.
func MouseIsPressed(b input.Button) bool { return mouseDown[b] }

// The current mouse position in cell coordinates.
func MousePosition() (x, y int) { return mouseX, mouseY }

// The mouse wheel scroll this tick. Up: 1, None: 0, Down: -1
func MouseScroll() int { return mouseScroll }

// private ========================================================

var mouseX, mouseY, mouseScroll int
var currentKeys []input.Key                      // key events this tick
var justPressedButtons = map[input.Button]bool{} // buttons that went down this tick
var mouseDown = map[input.Button]bool{}          // exact mouse button state

func clearInput() {
	clear(justPressedButtons)
	currentKeys = currentKeys[:0]
	mouseScroll = 0
}
func recordMouseButton(btn input.Button, down bool) {
	if down {
		if !mouseDown[btn] {
			justPressedButtons[btn] = true
		}
		mouseDown[btn] = true
	} else {
		delete(mouseDown, btn)
	}
}
func processEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		var mod = ev.Modifiers()

		if mod&tcell.ModShift != 0 {
			currentKeys = append(currentKeys, input.Shift)
		}
		if mod&tcell.ModCtrl != 0 {
			currentKeys = append(currentKeys, input.Control)
		}
		if mod&tcell.ModAlt != 0 {
			currentKeys = append(currentKeys, input.Alt)
		}

		if ev.Key() != tcell.KeyRune {
			currentKeys = append(currentKeys, input.Key(ev.Key()))
		} else if ev.Rune() != 0 {
			currentKeys = append(currentKeys, input.Key(ev.Rune()))
		}
	case *tcell.EventMouse:
		mouseX, mouseY = ev.Position()

		var btns = ev.Buttons()
		recordMouseButton(input.MouseLeft, btns&tcell.Button1 != 0)
		recordMouseButton(input.MouseRight, btns&tcell.Button2 != 0)
		recordMouseButton(input.MouseMiddle, btns&tcell.Button3 != 0)

		if btns&tcell.WheelUp != 0 {
			mouseScroll = 1
		}
		if btns&tcell.WheelDown != 0 {
			mouseScroll = 0
		}
	}
}
