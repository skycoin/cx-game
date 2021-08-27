package game

import (
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/particles/particle_emitter"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/sound"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/engine/ui/console"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/stars/starfield"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/world/worldgen"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

const (
	INITIAL_WINDOW_WIDTH  = 1440
	INITIAL_WINDOW_HEIGHT = 900
)

var (
	Console console.Console
	Cam     *camera.Camera
	win     render.Window
	window  *glfw.Window
	fps     *render.Fps
	player  *agents.Agent

	World              world.World
	DrawCollisionBoxes = false
	FPS                int

	//unused
	isTileSelectorVisible = false
	worldItem             *item.WorldItem
)

func Init() {
	win = render.NewWindow(INITIAL_WINDOW_WIDTH, INITIAL_WINDOW_HEIGHT, true)
	win.SetCallbacks()
	// defer glfw.Terminate()

	window = win.Window

	if runtime.GOOS == "darwin" {
		render.FixRenderCOCOA(window)
	}

	window.SetMouseButtonCallback(mouseButtonCallback)
	window.SetScrollCallback(scrollCallback)
	window.SetCursorPosCallback(cursorPosCallback)
	//window.SetSizeCallback(windowSizeCallback)

	input.Init(&win)
	sound.Init()
	sound.Mute()
	spriteloader.DEBUG = false
	spriteloader.InitSpriteloader(&win)
	spriteloader.LoadSpritesFromConfigs()
	anim.InitAnimatedSpriteLoader()
	world.Init()
	item.InitWorldItem()
	ui.Init()
	particles.InitParticles()
	item.RegisterItemTypes()
	render.Init()

	fps = render.NewFps(false)
	Cam = camera.NewCamera(&win)

	// TODO move this to the world package or similar
	World = worldgen.GenerateWorld()
	components.ChangeWorld(&World)

	//World.Planet = *mapgen.GeneratePlanet()
	Cam.PlanetWidth = float32(World.Planet.Width)

	starfield.InitStarField(&win, Cam)

	worldTiles := World.Planet.GetAllTilesUnique()
	log.Printf("Found [%v] unique tiles in the world", len(worldTiles))

	spawnX := int(20)
	Cam.SetCameraPosition(float32(spawnX), 5)

	spawnPos := cxmath.Vec2{
		float32(spawnX),
		float32(World.Planet.GetHeight(spawnX) + 10),
	}

	World.Entities.Agents.Spawn(
		constants.AGENT_TYPE_SLIME, agents.AgentCreationOptions{
			X: spawnPos.X - 6, Y: spawnPos.Y,
		},
	)
	World.Entities.Agents.Spawn(
		constants.AGENT_TYPE_SPIDER_DRILL, agents.AgentCreationOptions{
			X: spawnPos.X + 6, Y: spawnPos.Y,
		},
	)
	playerAgentID = World.Entities.Agents.Spawn(
		constants.AGENT_TYPE_PLAYER, agents.AgentCreationOptions{
			X: spawnPos.X, Y: spawnPos.Y,
		},
	)
	player = findPlayer()
	player.InventoryID = item.NewDevInventory()
	components.Init(&World, Cam, player)

	sound.LoadSound("player_jump", "jump.wav")
	Console = console.New()

	//add oxygen emitter
	particle_emitter.EmitOxygen(playerAgentID, &World.Entities.Particles)

}
