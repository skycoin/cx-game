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
	"github.com/skycoin/cx-game/particles"
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
	WINDOW_WIDTH  = 800
	WINDOW_HEIGHT = 480
)

var (
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

	isFreeCam = false
	//unused
	isTileSelectorVisible  = false
	isInventoryGridVisible = false
	tilePaletteSelector    ui.TilePaletteSelector

	worldItem *item.WorldItem

	inventoryId uint32
)

func main() {
	win = render.NewWindow(WINDOW_WIDTH, WINDOW_HEIGHT, true)
	defer glfw.Terminate()

	Init()
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetScrollCallback(scrollCallback)
	window.SetSizeCallback(windowSizeCallback)
	program := win.Program

	inventoryId = item.NewInventory(10, 8)

	debugItemType :=
		item.NewItemType(spriteloader.GetSpriteIdByName("Bedrock"))
	debugItemTypeId :=
		item.AddItemType(debugItemType)

	inventory := item.GetInventoryById(inventoryId)
	laserGunItemTypeId := item.RegisterLaserGunItemType()
	inventory.Slots[inventory.ItemSlotIndexForPosition(1, 7)] =
		item.InventorySlot{laserGunItemTypeId, 1, 0}


	worldTiles := CurrentPlanet.GetAllTilesUnique()
	log.Printf("Found [%v] unique tiles in the world", len(worldTiles))
	tilePaletteSelector = ui.
		NewDevTilePaleteSelector()

	//init cam and cat positions
	spawnX := int(20)
	Cam.X = float32(spawnX)
	Cam.Y = 5
	Cam.SetCameraPosition(Cam.X, Cam.Y)
	Cam.SetCameraZoomPosition(0)
	cat.Pos.X = float32(spawnX)
	cat.Pos.Y = float32(CurrentPlanet.GetHeight(spawnX) + 10)
	enemies.SpawnBasicEnemy(cat.Pos.X+6, cat.Pos.Y)
	enemies.SpawnBasicEnemy(cat.Pos.X-6, cat.Pos.Y)

	worldItem = item.NewWorldItem(debugItemTypeId)
	worldItem.Pos.X = cat.Pos.X - 3
	worldItem.Pos.Y = cat.Pos.Y + 2

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

func Init() {
	spriteloader.InitSpriteloader(&win)
	spriteloader.DEBUG = false
	item.InitWorldItem()
	ui.InitTextRendering()
	enemies.InitBasicEnemies()
	particles.InitParticles()

	cat = models.NewCat()
	fps = models.NewFps(false)
	Cam = camera.NewCamera(&win)
	CurrentPlanet = world.NewDevPlanet()

	starmap.Init(&win)
	starmap.Generate(256, 0.04, 8)
}

func Tick(dt float32) {
	Cam.Tick(dt)
	fps.Tick()
	if spacePressed {
		ui.PlaceDialogueBox(
			"*jump*", ui.AlignRight, 1,
			mgl32.Translate3D(
				cat.Pos.X,
				cat.Pos.Y,
				0,
			),
		)
	}
	if catIsScratching {
		ui.PlaceDialogueBox(
			"*scratch", ui.AlignLeft, 1,
			mgl32.Translate3D(
				cat.Pos.X,
				cat.Pos.Y,
				0,
			),
		)
	}
	ui.TickDialogueBoxes(dt)
	particles.TickParticles(dt)

	if worldItem != nil {
		pickupItem := worldItem.Tick(CurrentPlanet, dt, cat.Pos)
		if pickupItem {
			item.GetInventoryById(inventoryId).
				TryAddItem(worldItem.ItemTypeId)
			worldItem = nil
		}
	}
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
	enemies.TickBasicEnemies(CurrentPlanet, dt, cat, catIsScratching)

	spacePressed = false
	catIsScratching = false
}

func Draw(window *glfw.Window, program uint32, VAO uint32) {
	win.UpdateProjectionMatrix()
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	baseCtx := win.DefaultRenderContext()
	baseCtx.Projection = Cam.GetProjectionMatrix()
	//
	camCtx := baseCtx.PushView(Cam.GetView())

	starmap.Draw()
	CurrentPlanet.Draw(Cam,world.BgLayer)
	CurrentPlanet.Draw(Cam,world.MidLayer)
	// draw particles between mid and top layers.
	// this works quite well for weapon particle effects
	particles.DrawParticles(camCtx)
	CurrentPlanet.Draw(Cam,world.TopLayer)
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
	// FIXME: draw dialogue boxes uses alternate projection matrix;
	// restore original projection matrix
	win.UpdateProjectionMatrix()

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
			tilePaletteSelector.CycleLayer()
		}
		if k == glfw.KeyI {
			isInventoryGridVisible = !isInventoryGridVisible
		}
		if k == glfw.KeyLeftShift {
			catIsScratching = true
		}
		inventory := item.GetInventoryById(inventoryId)
		inventory.TrySelectSlot(k)
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
	if didSelectPaleteTile { return }

	if (tilePaletteSelector.IsMultiTileSelected()) {
		didPlaceMultiTile := CurrentPlanet.TryPlaceMultiTile(
			screenX, screenY,
			world.Layer(tilePaletteSelector.LayerIndex),
			tilePaletteSelector.GetSelectedMultiTile(),
			Cam,
		)
		if didPlaceMultiTile {
			return
		}
	} else {
		didPlaceTile := CurrentPlanet.TryPlaceTile(
			screenX, screenY,
			world.Layer(tilePaletteSelector.LayerIndex),
			tilePaletteSelector.GetSelectedTile(),
			Cam,
		)
		if didPlaceTile {
			return
		}
	}

	item.GetInventoryById(inventoryId).
		TryUseItem(screenX,screenY, Cam, CurrentPlanet, cat)
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

var yOffset float32 = 0

func scrollCallback(w *glfw.Window, xOff, yOff float64) {
	Cam.SetCameraZoomPosition(float32(yOff))
}

func boolToFloat(x bool) float32 {
	if x {
		return 1
	} else {
		return 0
	}
}
