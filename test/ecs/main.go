package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/cxecs"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

var (
	win            render.Window
	mouseX, mouseY float64
)

func init() {
	runtime.LockOSThread()
	rand.Seed(time.Now().Unix())
}

const (
	WIDTH  = 640 // 20 tiles wide
	HEIGHT = 480 // 15 tiles high
)

func main() {
	win = render.NewWindow(WIDTH, HEIGHT, false)
	window := win.Window
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)

	Init(win)

	dt := float32(0)
	lastFrame := float32(glfw.GetTime())

	// spriteId := spriteloader.GetSpriteIdByName("enemy")
	for !window.ShouldClose() {
		currentFrame := float32(glfw.GetTime())
		dt = currentFrame - lastFrame
		lastFrame = currentFrame

		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw()
		// spriteloader.DrawSpriteQuad(1, 1, 1, 1, spriteId)
		Update(dt)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func Init(win render.Window) {
	spriteloader.InitSpriteloader(&win)
	cxecs.InitDev(&win)
}

func Draw() {

}

func Update(dt float32) {
	//draws as well
	cxecs.UpdateDev(dt)
}

func cursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	mouseX = xpos
	mouseY = ypos
}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButton1 && action == glfw.Press {
		tileCoords := win.ConvertScreenToTileCoords(float32(mouseX), float32(mouseY))
		tileX := tileCoords.X()
		tileY := tileCoords.Y()
		fmt.Printf("Screen Coords: %v - %v    Tile Coords :  %v -%v\n", mouseX, mouseY, tileX, tileY)
	}
}
