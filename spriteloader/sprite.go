package spriteloader

type Sprite struct{
	name string
	xpos int
	ypos int
}

func NewSprite(name string, xpos,ypos int)*Sprite{
	return &Sprite{
		name: name,
		xpos: xpos,
		ypos: ypos,
	}
}