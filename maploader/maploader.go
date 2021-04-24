package main

import (
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

type Map struct {
	tiles []*MapTile
}

type MapTile struct {
	spriteId int
	x        int
	y        int
	show     int
}

var myMap Map

func InitMap() {
	myMap = Map{}

	spriteSheetId := spriteloader.LoadSpriteSheet("./maploader/planets.png")
	spriteloader.LoadSprite(spriteSheetId, "blue_star", 2, 2)
	// spriteloader.LoadSprite(spriteSheetId, "brown_star", 2, 1)
	// spriteloader.LoadSprite(spriteSheetId, "purple_star", 2, 0)
	// spriteloader.LoadSprite(spriteSheetId, "green_star", 2, 3)

	for i := 0; i < 16; i++ {
		show := 1
		if rand.Float32() < 0.15 {
			show = 0
		}
		myMap.tiles = append(myMap.tiles, &MapTile{
			spriteId: spriteloader.GetSpriteIdByName("blue_star"),
			x:        i % 4,
			y:        i / 4,
			show:     show,
		})
	}
}

func DrawMap() {
	// blueSpriteId := spriteloader.GetSpriteIdByName("blue_star")
	// brownSpriteId := spriteloader.GetSpriteIdByName("brown_star")
	// purpleSpriteId := spriteloader.GetSpriteIdByName("purple_star")
	// greenSpriteId := spriteloader.GetSpriteIdByName("green_star")

	// spriteloader.DrawSpriteQuad(-3, 3, 1, 1, blueSpriteId)
	// spriteloader.DrawSpriteQuad(-1, 3, 1, 1, brownSpriteId)
	// spriteloader.DrawSpriteQuad(1, 3, 1, 1, purpleSpriteId)
	// spriteloader.DrawSpriteQuad(3, 3, 1, 1, greenSpriteId)

	// field is 4x4 grid
	// -3, 3 is top left
	// 3,3 is top right
	// -3, -3 is bottom left
	// 3, -3 is bottom right

	//we need to convert coordinates 0,0 to -3, 3
	for _, tile := range myMap.tiles {
		xpos, ypos := convertCoords(tile.x, tile.y)
		if tile.show == 0 {
			continue
		}
		spriteloader.DrawSpriteQuad(xpos, ypos, 1, 1, tile.spriteId)
	}

}

func convertCoords(x int, y int) (int, int) {
	return -3 + (x * 2), 3 - (y * 2)
}

func main() {
	window := render.NewWindow(500, 500, false)
	defer glfw.Terminate()
	spriteloader.InitSpriteloader(&window)

	InitMap()

	// spriteSheetId := spriteloader.LoadSpriteSheet("./assets/starfield/stars/planets.png")
	// spriteloader.LoadSprite(spriteSheetId, "blue_star", 1, 1)

	// spriteId := spriteloader.GetSpriteIdByName("blue_star")

	for !window.Window.ShouldClose() {
		// gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// spriteloader.DrawSpriteQuad(-3, 3, 1, 1, spriteId)
		// spriteloader.DrawSpriteQuad(-1, 3, 1, 1, spriteId)

		// spriteloader.DrawSpriteQuad(1, -1, 1, 1, spriteId)

		DrawMap()
		glfw.PollEvents()
		window.Window.SwapBuffers()

	}
}
