package model

type Cat struct {
	name   string
	width  int
	height int
}

func NewCat() *Cat {
	cat := Cat{
		name:   "NewCat",
		width:  2,
		height: 2,
	}

	return &cat
}
