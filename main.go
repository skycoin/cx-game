package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/starmap"

	//cv "github.com/skycoin/cx-game/cmd/spritetool"

	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/ui"
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

var (
	dt, lastFrame float32
)

var upPressed bool
var downPressed bool
var leftPressed bool
var rightPressed bool
var spacePressed bool
var mouseX, mouseY float64

func mouseButtonCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	// we only care about mousedown events for now
	if a != glfw.Press {
		return
	}
	screenX := float32(2*mouseX/float64(win.Width) - 1)
	screenY := 1 - float32(2*mouseY/float64(win.Height))
	projection := win.GetProjectionMatrix()

	didSelectPaleteTile := tilePaletteSelector.
		TrySelectTile(screenX, screenY, projection)

	// only try to place a tile if we didn't select a palete with this click
	if !didSelectPaleteTile {
		CurrentPlanet.TryPlaceTile(
			screenX, screenY,
			projection,
			world.Layer(cyclingPalleteSelector),
			tilePaletteSelector.GetSelectedTile(),
			Cam,
		)
	}
}

func cursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	mouseX = xpos
	mouseY = ypos
}

var isFreeCam = false
var isTileSelectorVisible = false
var isInventoryGridVisible = false
var tilePaletteSelector ui.TilePaletteSelector
var cyclingPalleteSelector int = -1

var worldItem item.WorldItem
var cat *models.Cat
var fps *models.Fps

var Cam *camera.Camera
var win render.Window
var tex uint32

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press {
		if k == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
		if k == glfw.KeyW {
			upPressed = true
		}
		if k == glfw.KeyS {
			downPressed = true
		}
		if k == glfw.KeyA {
			leftPressed = true
		}
		if k == glfw.KeyD {
			rightPressed = true
		}
		if k == glfw.KeySpace {
			spacePressed = true
		}
		if k == glfw.KeyQ {
			Cam.Zoom += 0.5
		}
		if k == glfw.KeyZ {
			Cam.Zoom -= 0.5
		}
		if k == glfw.KeyF2 {
			isFreeCam = !isFreeCam
		}
		if k == glfw.KeyF3 {
			cyclingPalleteSelector++
			if cyclingPalleteSelector == 0 {
				tilePaletteSelector.Toggle()
			} else if cyclingPalleteSelector > 2 {
				cyclingPalleteSelector = -1
				tilePaletteSelector.Toggle()
			}
		}
		if k == glfw.KeyI {
			isInventoryGridVisible = !isInventoryGridVisible
		}
	} else if a == glfw.Release {
		if k == glfw.KeyW {
			upPressed = false
		}
		if k == glfw.KeyS {
			downPressed = false
		}
		if k == glfw.KeyA {
			leftPressed = false
		}
		if k == glfw.KeyD {
			rightPressed = false
		}
	}
}

var inventoryId uint32

func main() {

	/*
		var SS cv.SpriteSet
		SS.LoadFile("./assets/sprite.png", 250, false)
		SS.ProcessContours()
		SS.DrawSprite()
	*/

	win = render.NewWindow(width, height, true)
	spriteloader.InitSpriteloader(&win)
	ui.InitTextRendering()

	cat = models.NewCat()
	log.Printf("inventoryId=%v", inventoryId)
	fps = models.NewFps(false)

	CurrentPlanet = world.NewDevPlanet()
	inventoryId = item.NewInventory(10, 8)
	debugItemType :=
		item.NewItemType(spriteloader.GetSpriteIdByName("RedBlip"))

	inventory := item.GetInventoryById(inventoryId)
	inventory.Slots[inventory.ItemSlotIndexForPosition(0, 0)] =
		item.InventorySlot{debugItemType, 1}
	inventory.Slots[inventory.ItemSlotIndexForPosition(1, 7)] =
		item.InventorySlot{debugItemType, 5}

	worldTiles := CurrentPlanet.GetAllTilesUnique()
	log.Printf("Found [%v] unique tiles in the world", len(worldTiles))
	tilePaletteSelector = ui.
		MakeTilePaleteSelector(worldTiles)
	window := win.Window
	Cam = camera.NewCamera(&win)
	spawnX := int(20)
	Cam.X = float32(spawnX)
	Cam.Y = 5
	Cam.Zoom = -10
	cat.Pos.X = float32(spawnX)
	cat.Pos.Y = float32(CurrentPlanet.GetHeight(spawnX) + 10)

	worldItem = item.NewWorldItem(debugItemType)
	worldItem.Pos.X = cat.Pos.X - 3
	worldItem.Pos.Y = cat.Pos.Y + 2

	window.SetKeyCallback(keyCallBack)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	defer glfw.Terminate()
	// VAO := makeVao()
	program := win.Program
	gl.GenTextures(1, &tex)

	starmap.Init(&win)
	starmap.Generate(256, 0.04, 8)
	//without this panics
	lastFrame = float32(glfw.GetTime())

	for !window.ShouldClose() {
		currTime := float32(glfw.GetTime())
		dt = currTime - lastFrame
		lastFrame = currTime

		Tick()

		redraw(window, program, win.VAO)

		fps.Tick()

	}
}

func boolToFloat(x bool) float32 {
	if x {
		return 1
	} else {
		return 0
	}
}

func Tick() {
	worldItem.Tick(CurrentPlanet, dt)
	if isFreeCam {
		Cam.MoveCam(
			boolToFloat(rightPressed)-boolToFloat(leftPressed),
			boolToFloat(upPressed)-boolToFloat(downPressed),
			0,
			dt,
		)
		cat.Tick(false, false, false, CurrentPlanet, dt)
	} else {
		cat.Tick(leftPressed, rightPressed, spacePressed, CurrentPlanet, dt)
	}
	spacePressed = false
}

func redraw(window *glfw.Window, program uint32, VAO uint32) {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	fmt.Println(Cam.X, " ", Cam.Y, " ", Cam.Zoom)
	starmap.Draw()
	CurrentPlanet.Draw(Cam)
	worldItem.Draw(Cam)
	cat.Draw(Cam)

	// tile - air line (green)
	collidingTileLines := CurrentPlanet.GetCollidingTilesLinesRelative(
		int(cat.Pos.X), int(cat.Pos.Y))
	if len(collidingTileLines) > 2 {
		Cam.DrawLines(collidingTileLines, []float32{0.0, 1.0, 0.0})
	}

	// body bounding box (blue)
	Cam.DrawLines(cat.GetBBoxLines(), []float32{0.0, 0.0, 1.0})

	// colliding line from body (red)
	collidingLines := cat.GetCollidingLines()
	if len(collidingLines) > 2 {
		Cam.DrawLines(collidingLines, []float32{1.0, 0.0, 0.0})
	}

	inventory := item.GetInventoryById(inventoryId)
	if isInventoryGridVisible {
		inventory.DrawGrid()
	} else {
		inventory.DrawBar()
	}
	tilePaletteSelector.Draw()

	glfw.PollEvents()
	window.SwapBuffers()
}
