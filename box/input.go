package box

import (
	"retro-game-box/box/input"

	"github.com/gdamore/tcell/v2"
)

var currentInput []input.Input
var mouseX, mouseY int

func resetInput() {
	currentInput = currentInput[:0]
}

// processEvent records input from ev into the per-tick input state.
func processEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		var mod = ev.Modifiers()

		if mod&tcell.ModShift != 0 {
			currentInput = append(currentInput, input.Shift)
		}
		if mod&tcell.ModCtrl != 0 {
			currentInput = append(currentInput, input.Control)
		}
		if mod&tcell.ModAlt != 0 {
			currentInput = append(currentInput, input.Alt)
		}

		if ev.Key() != tcell.KeyRune {
			currentInput = append(currentInput, input.Input(ev.Key()))
		} else if ev.Rune() != 0 {
			currentInput = append(currentInput, input.Input(ev.Rune()))
		}
	case *tcell.EventMouse:
		mouseX, mouseY = ev.Position()

		var btns = ev.Buttons()
		if btns&tcell.Button1 != 0 {
			currentInput = append(currentInput, input.MouseLeft)
		}
		if btns&tcell.Button2 != 0 {
			currentInput = append(currentInput, input.MouseRight)
		}
		if btns&tcell.Button3 != 0 {
			currentInput = append(currentInput, input.MouseMiddle)
		}
		if btns&tcell.WheelUp != 0 {
			currentInput = append(currentInput, input.MouseWheelUp)
		}
		if btns&tcell.WheelDown != 0 {
			currentInput = append(currentInput, input.MouseWheelDown)
		}
	}
}

// Handles keyboard & mouse. Accepts runes.
func InputIsJustPressedAndHeld(input input.Input) bool {
	for _, i := range currentInput {
		if i == input {
			return true
		}
	}
	return false
}

func InputMousePosition() (x, y int) { return mouseX, mouseY }
