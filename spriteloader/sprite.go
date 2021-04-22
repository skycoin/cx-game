package spriteloader

import "image"

//Sprite is
type Sprite struct {
	name      string
	id        uint32
	xpos      int
	ypos      int
	imageInfo *image.RGBA
}
