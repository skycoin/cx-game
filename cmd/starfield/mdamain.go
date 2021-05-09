package main

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
)

func main() {
	window := render.NewWindow(800, 600, false)
	win := window.Window
	program := window.Program
	gl.UseProgram(program)
	spriteloader.InitSpriteloader(&window)

	ssId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/planets.png")

	spriteloader.LoadSprite(ssId, "abc", 1, 1)
	spriteloader.LoadSprite(ssId, "abcd", 2, 2)
	spriteId := spriteloader.GetSpriteIdByName("abc")
	config := make(map[string]interface{})
	mychan := make(chan struct{})

	go utility.CheckAndReload("./cmd/starfield/config/config.yaml", config, mychan)

	go func() {
		for {
			select {
			case <-mychan:
				id, ok := config["spriteid"]
				if !ok {
					log.Fatal("ERROR")
				}
				spriteId = id.(int)
			}
		}

	}()

	for !win.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT)

		spriteloader.DrawSpriteQuad(0, 0, 1, 1, spriteId)

		glfw.PollEvents()
		window.Window.SwapBuffers()
	}

}
