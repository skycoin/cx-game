package main

import (
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/render"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
)

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press && k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
}

func main() {
	log.Print("running test")
	win := render.NewWindow(640,480,true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	spriteSheetId := spriteloader.
		LoadSpriteSheet("./assets/starfield/stars/Starsheet1.png")
	spriteloader.
		LoadSprite(spriteSheetId, "star", 0,0)
	spriteId := spriteloader.
		GetSpriteIdByName("star")
	log.Print(spriteId)
	for !window.ShouldClose() {
		//DrawSpriteQuad(
		glfw.PollEvents()
	}
}
