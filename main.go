package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/input"
	"github.com/skycoin/cx-game/sound"
	"github.com/skycoin/cx-game/starfield"

	//cv "github.com/skycoin/cx-game/cmd/spritetool"

	"github.com/skycoin/cx-game/enemies"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/render/worldctx"
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
	WINDOW_HEIGHT = 600
)

var (
	Cam    *camera.Camera
	win    render.Window
	player *models.Player
	fps    *models.Fps

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

	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetScrollCallback(scrollCallback)
	window.SetSizeCallback(windowSizeCallback)

	inventoryId = item.NewDevInventory()

	worldTiles := CurrentPlanet.GetAllTilesUnique()
	log.Printf("Found [%v] unique tiles in the world", len(worldTiles))
	tilePaletteSelector = ui.
		NewDevTilePaleteSelector()

	//init cam and cat positions
	spawnX := int(20)
	Cam.SetCameraPosition(float32(spawnX), 5)
	// Cam.SetCameraZoomPosition(0)
	player.Pos.X = float32(spawnX)
	player.Pos.Y = float32(CurrentPlanet.GetHeight(spawnX) + 10)
	enemies.SpawnBasicEnemy(player.Pos.X+6, player.Pos.Y)
	enemies.SpawnBasicEnemy(player.Pos.X-6, player.Pos.Y)

	sound.LoadSound("player_jump", "jump.wav")

	var dt, lastFrame float32
	lastFrame = float32(glfw.GetTime())

	for !window.ShouldClose() {
		currTime := float32(glfw.GetTime())
		dt = currTime - lastFrame
		lastFrame = currTime

		ProcessInput()
		Tick(dt)
		Draw()
	}
}

func Init() {
	input.Init(win.Window)
	sound.Init()
	spriteloader.InitSpriteloader(&win)
	spriteloader.DEBUG = false
	item.InitWorldItem()
	ui.InitTextRendering()
	enemies.InitBasicEnemies()
	particles.InitParticles()
	item.RegisterItemTypes()

	player = models.NewPlayer()

	fps = models.NewFps(false)
	Cam = camera.NewCamera(&win)
	//CurrentPlanet = world.NewDevPlanet()
	CurrentPlanet = world.GeneratePlanet()
	Cam.PlanetWidth = float32(CurrentPlanet.Width)

	starfield.InitStarField(&win, player, Cam)

}

func FixedTick() {
	if Cam.IsFreeCam() {
		player.Controlled = false
	} else {
		player.Controlled = true
	}
	player.FixedTick(CurrentPlanet)

}

func Tick(dt float32) {
	// TODO account for the scenario where multiple physics ticks are required
	if physics.WillTick(dt) {
		FixedTick()
	}
	physics.Simulate(dt, CurrentPlanet)
	if Cam.IsFreeCam() {
		Cam.MoveCam(dt)
	} else {
		playerPos := player.InterpolatedTransform.Col(3).Vec2()
		Cam.SetCameraPosition(playerPos.X(), playerPos.Y())
	}
	Cam.Tick(dt)
	fps.Tick()
	ui.TickDialogueBoxes(dt)
	particles.TickParticles(dt)
	pickedUpItems := item.TickWorldItems(CurrentPlanet, dt, player.Pos)
	for _, worldItem := range pickedUpItems {
		item.GetInventoryById(inventoryId).TryAddItem(worldItem.ItemTypeId)
	}
	enemies.TickBasicEnemies(CurrentPlanet, dt, player, catIsScratching)

	sound.SetListenerPosition(player.Pos)
	//has to be after listener position is updated
	sound.Update()

	starfield.UpdateStarField(dt)
	catIsScratching = false

}

func Draw() {
	gl.ClearColor(7.0/255.0, 8.0/255.0, 25.0/255.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	baseCtx := win.DefaultRenderContext()
	baseCtx.Projection = Cam.GetProjectionMatrix()
	camCtx := baseCtx.PushView(Cam.GetView())
	worldCtx := worldctx.NewWorldRenderContext(Cam, CurrentPlanet)

	starfield.DrawStarField()
	CurrentPlanet.Draw(Cam, world.BgLayer)
	CurrentPlanet.Draw(Cam, world.MidLayer)
	// draw lasers between mid and top layers.
	particles.DrawMidTopParticles(worldCtx)
	CurrentPlanet.Draw(Cam, world.TopLayer)
	particles.DrawTopParticles(camCtx)
	/*
		if worldItem != nil {
			worldItem.Draw(Cam)
		}
	*/
	item.DrawWorldItems(Cam)
	enemies.DrawBasicEnemies(Cam)
	player.Draw(Cam, CurrentPlanet)

	// tile - air line (green)
	collidingTileLines := CurrentPlanet.GetCollidingTilesLinesRelative(
		int(player.Pos.X), int(player.Pos.Y))
	if len(collidingTileLines) > 2 {
		Cam.DrawLines(collidingTileLines, []float32{0.0, 1.0, 0.0}, baseCtx)
	}

	// body bounding box (blue)
	Cam.DrawLines(player.GetBBoxLines(), []float32{0.0, 0.0, 1.0}, baseCtx)

	// colliding line from body (red)
	collidingLines := player.GetCollidingLines()
	if len(collidingLines) > 2 {
		Cam.DrawLines(collidingLines, []float32{1.0, 0.0, 0.0}, baseCtx)
	}

	ui.DrawDialogueBoxes(camCtx)
	// FIXME: draw dialogue boxes uses alternate projection matrix;
	// restore original projection matrix

	inventory := item.GetInventoryById(inventoryId)
	if isInventoryGridVisible {
		inventory.DrawGrid(win.DefaultRenderContext())
	} else {
		inventory.DrawBar(win.DefaultRenderContext())
	}
	tilePaletteSelector.Draw(win.DefaultRenderContext())

	glfw.PollEvents()
	win.Window.SwapBuffers()
}

func ProcessInput() {
	if input.GetButtonDown("switch-helmet") {
		player.SetHelmNext()
	}
	if input.GetButtonDown("switch-suit") {
		player.SetSuitNext()
	}
	if input.GetButtonDown("jump") {
		didJump := player.Jump()
		if didJump {
			ui.PlaceDialogueBox(
				"*jump*", ui.AlignRight, 1,
				mgl32.Translate3D(
					player.Pos.X,
					player.Pos.Y,
					0,
				),
			)
			sound.PlaySound("player_jump", sound.SoundOptions{Pitch: 1.5})
		}
	}
	if input.GetButtonDown("fly") {
		player.ToggleFlying()
	}

	if input.GetButtonDown("scratch") {
		ui.PlaceDialogueBox(
			"*scratch", ui.AlignLeft, 1,
			mgl32.Translate3D(
				player.Pos.X,
				player.Pos.Y,
				0,
			),
		)
		catIsScratching = true
	}
	if input.GetButtonDown("mute") {
		sound.ToggleMute()
	}
	if input.GetButtonDown("freecam") {
		Cam.ToggleFreeCam()
	}
	if input.GetButtonDown("cycle-palette") {
		tilePaletteSelector.CycleLayer()
	}
	if input.GetButtonDown("inventory-grid") {
		isInventoryGridVisible = !isInventoryGridVisible
	}
	if input.GetKeyDown(glfw.KeyL) {
		starfield.SwitchBackgrounds(starfield.BACKGROUND_NEBULA)
	}
	if input.GetKeyDown(glfw.KeyO) {
		starfield.SwitchBackgrounds(starfield.BACKGROUND_VOID)
	}
	inventory := item.GetInventoryById(inventoryId)
	inventory.TrySelectSlot(input.GetLastKey())
}

type mouseDraws struct {
	xpos float32
	ypos float32
}

func mouseButtonCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	if a == glfw.Press {
		mousePressCallback(w, b, a, mk)
	}
	if a == glfw.Release {
		mouseReleaseCallback(w, b, a, mk)
	}
}

func mouseReleaseCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	screenX := float32(input.GetMouseX()-float64(win.Width)/2) / Cam.Zoom // adjust mouse position with zoom
	screenY := (float32(input.GetMouseY()-float64(win.Height)/2) * -1) / Cam.Zoom

	if isInventoryGridVisible {
		inventory := item.GetInventoryById(inventoryId)
		inventory.TryMoveSlot(screenX, screenY, Cam, CurrentPlanet, player)
	}
}

func mousePressCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	// we only care about mousedown events for now
	if a != glfw.Press {
		return
	}

	screenX := float32(((input.GetMouseX()-float64(widthOffset))/float64(scale) - float64(win.Width)/2)) / Cam.Zoom // adjust mouse position with zoom
	screenY := float32(((input.GetMouseY()-float64(heightOffset))/float64(scale)-float64(win.Height)/2)*-1) / Cam.Zoom

	didSelectPaleteTile := tilePaletteSelector.TrySelectTile(screenX, screenY)
	if didSelectPaleteTile {
		return
	}

	if tilePaletteSelector.IsMultiTileSelected() {
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

	if isInventoryGridVisible {
		inventory := item.GetInventoryById(inventoryId)
		clickedSlot :=
			inventory.TryClickSlot(screenX, screenY, Cam, CurrentPlanet, player)
		if clickedSlot {
			return
		}
	}

	item.GetInventoryById(inventoryId).
		TryUseItem(screenX, screenY, Cam, CurrentPlanet, player)
}

var (
	widthOffset, heightOffset int32
	scale                     float32 = 1
)

func windowSizeCallback(window *glfw.Window, width, height int) {

	// gl.Viewport(0, 0, int32(width), int32(height))
	scaleToFitWidth := float32(width) / float32(win.Width)
	scaleToFitHeight := float32(height) / float32(win.Height)
	scale = cxmath.Min(scaleToFitHeight, scaleToFitWidth)

	widthOffset = int32((float32(width) - float32(win.Width)*scale) / 2)
	heightOffset = int32((float32(height) - float32(win.Height)*scale) / 2)
	//correct mouse offsets
	input.UpdateMouseCoords(widthOffset, heightOffset, scale)

	gl.Viewport(widthOffset, heightOffset, int32(float32(win.Width)*scale), int32(float32(win.Height)*scale))
	// win.Width = width
	// win.Height = height
}

func scrollCallback(w *glfw.Window, xOff, yOff float64) {
	Cam.SetCameraZoomPosition(float32(yOff))
}
