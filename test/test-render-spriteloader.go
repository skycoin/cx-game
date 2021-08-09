package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
)

var spriteIndex int = 1

func init() {
	runtime.LockOSThread()
}

func keyCallBack(
	w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey,
) {
	if a != glfw.Press {
		return
	}
	if k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	if k == glfw.KeySpace {
		spriteIndex++
	}
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	win := render.NewWindow(600, 400, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	spriteloader.InitSpriteloader(&win)
	render.Init()
	spritesheet := spriteloader.LoadSpriteSheetFromConfig(
		"./assets/sprite/containers.png", "./assets/sprite/containers.yaml",
	)
	for !window.ShouldClose() {
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// bind texture
		gl.Enable(gl.TEXTURE_2D)
		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		// NOTE depth test is disabled for now,
		// as we assume that objects are drawn in the correct order
		gl.Disable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)

		bucket, ok := spritesheet.Sprite("orangebucket")
		if !ok {
			log.Fatal("could not find orangebucket sprite")
		}
		ac, ok := spritesheet.Sprite("ac-unit")
		if !ok {
			log.Fatal("could not find ac-unit sprite")
		}

		spritesheet.Draw([]render.SpriteRenderParams{
			render.SpriteRenderParams{
				MVP:    win.DefaultRenderContext().MVP(),
				Sprite: bucket,
			},
			render.SpriteRenderParams{
				MVP: win.DefaultRenderContext().
					PushLocal(mgl32.Translate3D(2, 0, 0)).MVP(),
				Sprite: ac,
			},
		})

		glfw.PollEvents()
		window.SwapBuffers()
	}

}
