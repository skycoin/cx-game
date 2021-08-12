package game

import "github.com/go-gl/glfw/v3.3/glfw"

func Run() {
	Init()

	dt := float32(0)
	lastFrame := float32(glfw.GetTime())

	for !window.ShouldClose() {
		currTime := float32(glfw.GetTime())
		dt = currTime - lastFrame
		lastFrame = currTime

		ProcessInput()
		Update(dt)
		Draw()
	}
}
