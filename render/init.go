package render

func Init() {
	initSprite()
	initColor()

	InitMainFramebuffer()
	InitPlanetFrameBuffer()
	InitOutlineProgram()
}
