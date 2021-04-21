package model

type Screen struct {
	x      int
	y      int
	height int
	width  int
}

func NewScreen() *Screen {
	screen := Screen{
		x:      0,
		y:      0,
		height: 480,
		width:  800,
	}

	return &screen
}
