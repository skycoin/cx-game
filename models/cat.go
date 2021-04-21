package model

type Cat struct {
	Name   string
	Width  int
	Height int

	XVelocity float32
	YVelocity float32
}

func NewCat() *Cat {
	cat := Cat{
		Name:      "NewCat",
		Width:     2,
		Height:    2,
		XVelocity: 0.0,
		YVelocity: 0.0,
	}

	return &cat
}
