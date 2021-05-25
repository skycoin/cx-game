package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/starmap"

	//cv "github.com/skycoin/cx-game/cmd/spritetool"

	"github.com/skycoin/cx-game/enemies"
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

const (
	width  = 800
	height = 480
)

var (
	dt, lastFrame float32

	Cam *camera.Camera
	win render.Window
	cat *models.Cat
	fps *models.Fps

	CurrentPlanet      *world.Planet
	DrawCollisionBoxes = false
	FPS                int

	catIsScratching bool

	upPressed      bool
	downPressed    bool
	leftPressed    bool
	rightPressed   bool
	spacePressed   bool
	mouseX, mouseY float64

	isFreeCam              = false
	isTileSelectorVisible  = false
	isInventoryGridVisible = false
	tilePaletteSelector    ui.TilePaletteSelector
	cyclingPalleteSelector int = -1

	worldItem *item.WorldItem

	inventoryId uint32
)

func main() {

	/*
		var SS cv.SpriteSet
		SS.LoadFile("./assets/sprite.png", 250, false)
		SS.ProcessContours()
		SS.DrawSprite()
	*/

	win = render.NewWindow(width, height, true)
	defer glfw.Terminate()
	spriteloader.InitSpriteloader(&win)
	spriteloader.DEBUG = false

	item.InitWorldItem()
	ui.InitTextRendering()
	enemies.InitBasicEnemies()

	cat = models.NewCat()
	fps = models.NewFps(false)

	window := win.Window
	window.SetKeyCallback(keyCallBack)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetScrollCallback(scrollCallback)

	program := win.Program
	Cam = camera.NewCamera(&win)
	CurrentPlanet = world.NewDevPlanet()

	inventoryId = item.NewInventory(10, 8)
	debugItemType :=
		item.NewItemType(spriteloader.GetSpriteIdByName("RedBlip"))

	inventory := item.GetInventoryById(inventoryId)
	inventory.Slots[inventory.ItemSlotIndexForPosition(1, 7)] =
		item.InventorySlot{uint32(debugItemType), 5}

	worldTiles := CurrentPlanet.GetAllTilesUnique()
	log.Printf("Found [%v] unique tiles in the world", len(worldTiles))
	tilePaletteSelector = ui.
		MakeTilePaleteSelector(worldTiles)

	spawnX := int(20)
	Cam.X = float32(spawnX)
	Cam.Y = 5
	Cam.SetCameraPosition(Cam.X, Cam.Y)
	// Cam.Zoom = -10
	cat.Pos.X = float32(spawnX)
	cat.Pos.Y = float32(CurrentPlanet.GetHeight(spawnX) + 10)

	enemies.SpawnBasicEnemy(cat.Pos.X+6, cat.Pos.Y)
	enemies.SpawnBasicEnemy(cat.Pos.X-6, cat.Pos.Y)

	worldItem = item.NewWorldItem(uint32(debugItemType))
	worldItem.Pos.X = cat.Pos.X - 3
	worldItem.Pos.Y = cat.Pos.Y + 2

	window.SetKeyCallback(keyCallBack)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetSizeCallback(windowSizeCallback)
	defer glfw.Terminate()
	// VAO := makeVao()
	//gl.GenTextures(1, &tex)

	starmap.Init(&win)
	starmap.Generate(256, 0.04, 8)
	//without this panics

	var dt, lastFrame float32
	lastFrame = float32(glfw.GetTime())

	for !window.ShouldClose() {
		currTime := float32(glfw.GetTime())
		dt = currTime - lastFrame
		lastFrame = currTime

		Tick(dt)

		Draw(window, program, win.VAO)

	}
}

func boolToFloat(x bool) float32 {
	if x {
		return 1
	} else {
		return 0
	}
}

func Tick(dt float32) {

	if spacePressed {
		ui.PlaceDialogueBox(
			"*jump*", ui.AlignRight, 1,
			mgl32.Translate3D(
				cat.Pos.X,
				cat.Pos.Y,
				-spriteloader.SpriteRenderDistance,
			),
		)
	}
	if catIsScratching {
		ui.PlaceDialogueBox(
			"*scratch", ui.AlignLeft, 1,
			mgl32.Translate3D(
				cat.Pos.X,
				cat.Pos.Y,
				-spriteloader.SpriteRenderDistance,
			),
		)
	}
	ui.TickDialogueBoxes(dt)

	if worldItem != nil {
		pickupItem := worldItem.Tick(CurrentPlanet, dt, cat.Pos)
		if pickupItem {
			item.GetInventoryById(inventoryId).
				TryAddItem(worldItem.ItemTypeId)
			worldItem = nil
		}
	}

	enemies.TickBasicEnemies(CurrentPlanet, dt, cat, catIsScratching)

	if isFreeCam {
		Cam.MoveCam(
			boolToFloat(rightPressed)-boolToFloat(leftPressed),
			boolToFloat(upPressed)-boolToFloat(downPressed),
			dt,
		)
		cat.Tick(false, false, false, CurrentPlanet, dt)
	} else {
		cat.Tick(leftPressed, rightPressed, spacePressed, CurrentPlanet, dt)
	}
	spacePressed = false

	fps.Tick()
	catIsScratching = false

}

func Draw(window *glfw.Window, program uint32, VAO uint32) {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// fmt.Println(Cam.X, " ", Cam.Y, " ", Cam.Zoom)
	baseCtx := win.DefaultRenderContext()
	camCtx := baseCtx.PushView(Cam.GetView())
	//fmt.Println(Cam.X, " ", Cam.Y, " ", Cam.Zoom)
	starmap.Draw()
	CurrentPlanet.Draw(Cam)
	if worldItem != nil {
		worldItem.Draw(Cam)
	}
	enemies.DrawBasicEnemies(Cam)
	cat.Draw(Cam)

	// tile - air line (green)
	collidingTileLines := CurrentPlanet.GetCollidingTilesLinesRelative(
		int(cat.Pos.X), int(cat.Pos.Y))
	if len(collidingTileLines) > 2 {
		Cam.DrawLines(collidingTileLines, []float32{0.0, 1.0, 0.0}, baseCtx)
	}

	// body bounding box (blue)
	Cam.DrawLines(cat.GetBBoxLines(), []float32{0.0, 0.0, 1.0}, baseCtx)

	// colliding line from body (red)
	collidingLines := cat.GetCollidingLines()
	if len(collidingLines) > 2 {
		Cam.DrawLines(collidingLines, []float32{1.0, 0.0, 0.0}, baseCtx)
	}

	ui.DrawDialogueBoxes(camCtx)
	inventory := item.GetInventoryById(inventoryId)
	if isInventoryGridVisible {
		inventory.DrawGrid(baseCtx)
	} else {
		inventory.DrawBar(baseCtx)
	}
	tilePaletteSelector.Draw(baseCtx)

	glfw.PollEvents()
	window.SwapBuffers()
}

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
		if k == glfw.KeyLeftShift {
			catIsScratching = true
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

func mouseButtonCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	// we only care about mousedown events for now
	if a != glfw.Press {
		return
	}

	screenX := float32(mouseX - float64(win.Width)/2)
	screenY := float32(mouseY-float64(win.Height)/2) * -1

	didSelectPaleteTile := tilePaletteSelector.TrySelectTile(screenX, screenY)

	// only try to place a tile if we didn't select a palete with this click
	if !didSelectPaleteTile {
		CurrentPlanet.TryPlaceTile(
			screenX, screenY,
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

func windowSizeCallback(window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
	win.Width = width
	win.Height = height
}

func scrollCallback(w *glfw.Window, xpos, ypos float64) {
	Cam.SetCameraZoomPosition(float32(ypos))
}
