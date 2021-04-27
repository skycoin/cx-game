package main

import (
	"log"
	"runtime"

	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called 
	// from the main thread.
	runtime.LockOSThread()
}

var earth *world.Planet
var cam *camera.Camera

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press && k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
}

func main() {
	log.Print("running test")
	log.Print("You should see 3 overlapping tiles from different layers.")
	log.Print("top-to-bottom: (blue, orange, pink)")
	win := render.NewWindow(640,480,true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	cam = camera.NewCamera()
	cam.X = 2
	cam.Y = 2
	earth = world.NewPlanet(4,4)

	spriteloader.InitSpriteloader(&win)
	spriteSheetId := spriteloader.
		LoadSpriteSheet("./assets/starfield/stars/planets.png")
	spriteloader.
		LoadSprite(spriteSheetId, "big", 2,0)
	bigSpriteId := spriteloader.
		GetSpriteIdByName("big")
	spriteloader.
		LoadSprite(spriteSheetId, "mid", 1,1)
	midSpriteId := spriteloader.
		GetSpriteIdByName("mid")
	spriteloader.
		LoadSprite(spriteSheetId, "small", 3,2)
	smallSpriteId := spriteloader.
		GetSpriteIdByName("small")

	log.Print(smallSpriteId,midSpriteId,bigSpriteId)

	earth.Layers.Top[0] = world.Tile{
		SpriteID: uint32(smallSpriteId),
		TileType: world.TileTypeNormal,
	}
	earth.Layers.Mid[0]= world.Tile{
		SpriteID: uint32(midSpriteId),
		TileType: world.TileTypeNormal,
	}
	earth.Layers.Background[0] = world.Tile{
		SpriteID: uint32(bigSpriteId),
		TileType: world.TileTypeNormal,
	}

	for !window.ShouldClose() {
		gl.ClearColor(1,1,1,1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		earth.Draw(cam)
		glfw.PollEvents()
		window.SwapBuffers()
	}
}
