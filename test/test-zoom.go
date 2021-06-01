package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/world"
)

var (
	myPlanet   *world.Planet
	Cam        *camera.Camera
	shader     *utility.Shader
	width              = 800
	height             = 600
	zoom       float64 = 1
	projection mgl32.Mat4
)

func main() {
	window := render.NewWindow(width, height, false)
	glfwWindow := window.Window

	glfwWindow.SetScrollCallback(scrollCallback)
	glfwWindow.SetKeyCallback(keyCallback)
	spriteloader.InitSpriteloader(&window)

	shader = utility.ShaderWrapper(window.Program)
	shader.Use()
	Cam = camera.NewCamera(&window)
	Cam.X = 15
	Cam.Y = 15

	// projection = window.DefaultRenderContext().Projection
	projection = mgl32.Ortho(-float32(width)/32/2, float32(width)/32/2, -float32(height)/32/2, float32(height)/32/2, -1, 1000)
	shader.SetMat4("projection", &projection)

	spritesheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/planets.png")
	spriteloader.LoadSprite(spritesheetId, "tile", 1, 1)
	spriteId := spriteloader.GetSpriteIdByName("tile")
	for !glfwWindow.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		for y := 0; y < 2; y++ {
			for x := 0; x < 2; x++ {
				spriteloader.DrawSpriteQuad(float32(x), float32(y), 1, 1, spriteId)
			}
		}

		glfw.PollEvents()
		glfwWindow.SwapBuffers()
	}
}

func scrollCallback(w *glfw.Window, xpos, ypos float64) {
	zoom += ypos * 0.05
	fmt.Println(projection)
	ortho := mgl32.Ortho(-float32(width)/32/2/float32(zoom), float32(width)/32/2/float32(zoom), -float32(height)/32/2/float32(zoom), float32(height)/32/2/float32(zoom), -1, 1000)
	// ortho := projection.Mul(float32(zoom))
	shader.SetMat4("projection", &ortho)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

}
