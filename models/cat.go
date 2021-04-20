package model

type Cat struct {
	name string
	x    int
	y    int
}

func NewCat() *Cat {
	cat := Cat{
		name: "",
		x:    2,
		y:    2,
	}

	return &cat
}
