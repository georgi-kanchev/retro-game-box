package input

type Input int16

const (
	None Input = iota

	MouseLeft
	MouseRight
	MouseMiddle
	MouseWheelUp
	MouseWheelDown
)
const (
	Backspace = 8
	Tab       = 9
	Enter     = 13
	Escape    = 27
	Shift     = 28
	Control   = 29
	Alt       = 30
)

const (
	Space Input = iota + 32
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
	UpArrow Input = iota + 257
	DownArrow
	RightArrow
	LeftArrow
)
const (
	PageUp Input = iota + 266
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
	Menu Input = iota + 343
	CapsLock
	ScrollLock
	NumLock
)
