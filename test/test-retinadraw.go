package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called
	// from the main thread.
	runtime.LockOSThread()
}

const (
	SCREEN_HEIGHT = 640
	SCREEN_WIDTH  = 480
)

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	window, _ := glfw.CreateWindow(int(800), int(600), "Test RETINA DRAW", nil, nil)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}

	actualScreenWidth, actualScreenHeight := glfw.GetCurrentContext().GetFramebufferSize()
	gl.Viewport(0, 0, int32(actualScreenHeight), int32(actualScreenWidth))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, SCREEN_HEIGHT, 0, SCREEN_WIDTH, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	var vertices = []float32{
		200, 300, 0.0,
		300, 300, 0.0,
		200, 250, 0.0,
		300, 250, 0.0,

		320, 150, 0.0,
		350, 200, 0.0,

		400, 50, 0.0,
		450, 150, 0.0,
	}

	// font, err := glfont.LoadFont("Monofur for Powerline.ttf", int32(52), 800, 600)
	// if err != nil {
	// 	log.Panicf("LoadFont: %v", err)
	// }

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.EnableClientState(gl.VERTEX_ARRAY)
		gl.VertexPointer(3, gl.FLOAT, 0, gl.Ptr(vertices))
		gl.DrawArrays(gl.QUAD_STRIP, 0, 8)
		gl.DisableClientState(gl.VERTEX_ARRAY)
		// font.SetColor(1.0, 1.0, 1.0, 1.0)
		// font.Printf(100, 100, 1.0, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")

		window.SwapBuffers()
		glfw.PollEvents()

	}
}
