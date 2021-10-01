package render

var ScreenShader Program

func Init() {
	initSprite()
	initColor()

	screenShaderConfig := NewShaderConfig("./assets/shader/postprocessing.vert", "./assets/shader/postprocessing.frag")
	ScreenShader = screenShaderConfig.Compile()

	InitMainFramebuffer()
	InitPlanetFrameBuffer()
	InitOutlineProgram()
}
