package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var (
	DrawCollisionBoxes = false
	FPS                int
)

var CurrentPlanet *world.Planet

const (
	width  = 800
	height = 480
)

func main() {

	win := render.NewWindow(height, width, true)
	window := win.Window
	defer glfw.Terminate()

	program := win.Program

	for !window.ShouldClose() {
		Tick()
		draw(window, program)
	}
}

func Tick() {

}

func draw(window *glfw.Window, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	imgFile, err := os.Open("./assets/Daco_3555790.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()

	// Decode detexts the type of image as long as its image/<type> is imported
	im, err := png.Decode(imgFile)

	img := im.(*image.NRGBA)

	if err != nil {
		log.Fatalln(err)
	}

	size := img.Rect.Size()

	// flip image: first pixel is lower left corner
	data := make([]byte, size.X*size.Y*4)
	lineLen := size.X * 4
	dest := len(data) - lineLen
	for src := 0; src < len(img.Pix); src += img.Stride {
		copy(data[dest:dest+lineLen], img.Pix[src:src+img.Stride])
		dest -= lineLen
	}

	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))

	glfw.PollEvents()
	window.SwapBuffers()
}
