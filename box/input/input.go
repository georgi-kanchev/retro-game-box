package input

type Button uint16 // Mouse button.
type Key uint16    // Keyboard key.

const (
	None Button = iota

	MouseLeft
	MouseRight
	MouseMiddle
)
const (
	Backspace Key = 8
	Tab       Key = 9
	Enter     Key = 13
	Escape    Key = 27
	Shift     Key = 28
	Control   Key = 29
	Alt       Key = 30
)

const (
	Space Key = iota + 32
	ExclamationMark
	DoubleQuote
	Hash
	Dollar
	Percent
	Ampersand
	Quote
	LeftBrace
	RightBrace
	Asterisk
	Plus
	Comma
	Minus
	Dot
	Slash
	Number0
	Number1
	Number2
	Number3
	Number4
	Number5
	Number6
	Number7
	Number8
	Number9
	Colon
	Semicolon
	LessThan
	Equal
	MoreThan
	QuestionMark
	At
	UppercaseA
	UppercaseB
	UppercaseC
	UppercaseD
	UppercaseE
	UppercaseF
	UppercaseG
	UppercaseH
	UppercaseI
	UppercaseJ
	UppercaseK
	UppercaseL
	UppercaseM
	UppercaseN
	UppercaseO
	UppercaseP
	UppercaseQ
	UppercaseR
	UppercaseS
	UppercaseT
	UppercaseU
	UppercaseV
	UppercaseW
	UppercaseX
	UppercaseY
	UppercaseZ
	LeftSquareBrace
	Backslash
	RightSquareBrace
	Carat
	Underscore
	Backtick
	LowercaseA
	LowercaseB
	LowercaseC
	LowercaseD
	LowercaseE
	LowercaseF
	LowercaseG
	LowercaseH
	LowercaseI
	LowercaseJ
	LowercaseK
	LowercaseL
	LowercaseM
	LowercaseN
	LowercaseO
	LowercaseP
	LowercaseQ
	LowercaseR
	LowercaseS
	LowercaseT
	LowercaseU
	LowercaseV
	LowercaseW
	LowercaseX
	LowercaseY
	LowercaseZ
	LeftCurlyBrace
	Pipe
	RightCurlyBrace
	Tilde
)
const (
	UpArrow Key = iota + 257
	DownArrow
	RightArrow
	LeftArrow
)
const (
	PageUp Key = iota + 266
	PageDn
	Home
	End
	Insert
	Delete
	Help
	Exit
	Clear
	Cancel
	Print
	Pause
	Backtab
	F1
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
)
const (
	Menu Key = iota + 343
	CapsLock
	ScrollLock
	NumLock
)
