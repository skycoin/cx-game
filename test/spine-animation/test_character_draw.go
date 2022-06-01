package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	// OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/engine/sound"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/engine/ui/console"
	"github.com/skycoin/cx-game/game"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/stars/starfield"
	"github.com/skycoin/cx-game/test/spine-animation/animation"
	c "github.com/skycoin/cx-game/test/spine-animation/character"
	"github.com/skycoin/cx-game/world"
	worldimport "github.com/skycoin/cx-game/world/import"
)

const (
	width  = 500
	height = 500
)

var (
	characters     []*c.Character
	character      *c.Character
	characterIndex int
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
	// isTileSelectorVisible = false
	// worldItem             *item.WorldItem
)

func init1() {
	log.Printf("here")
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called
	// from the main thread.
	//runtime.LockOSThread()

	flags := game.ParseStartupFlags()

	vvw := int(constants.VIRTUAL_VIEWPORT_WIDTH)
	vvh := int(constants.VIRTUAL_VIEWPORT_HEIGHT)
	win = render.NewWindow(vvw, vvh, true)
	win.SetCallbacks()
	// defer glfw.Terminate()

	window = win.Window

	if runtime.GOOS == "darwin" {
		render.FixRenderCOCOA(window)

	}

	// window.SetMouseButtonCallback(mouseButtonCallback)
	// window.SetScrollCallback(scrollCallback)
	// window.SetCursorPosCallback(cursorPosCallback)
	//window.SetSizeCallback(windowSizeCallback)

	// input.Init(win.Window)
	// sound.Init()
	// sound.Mute()
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
	//World = worldgen.GenerateWorld()
	World = worldimport.ImportWorld(flags.TmxPath)
	World.Planet.InitLighting()
	World.Planet.DetectCircuits()
	components.ChangeWorld(&World)

	//World.Planet = *mapgen.GeneratePlanet()
	Cam.PlanetWidth = float32(World.Planet.Width)

	starfield.InitStarField(&win, Cam)

	worldTiles := World.Planet.GetAllTilesUnique()
	log.Printf("Found [%v] unique tiles in the world", len(worldTiles))

	spawnPos := cxmath.Vec2{80, 109} // start pos for moon bunker map

	World.Entities.Agents.Spawn(
		constants.AGENT_TYPE_SLIME, agents.AgentCreationOptions{
			X: spawnPos.X - 10, Y: spawnPos.Y,
		},
	)
	World.Entities.Agents.Spawn(
		constants.AGENT_TYPE_SPIDER_DRILL, agents.AgentCreationOptions{
			X: spawnPos.X + 6, Y: spawnPos.Y,
		},
	)
	World.Entities.Agents.Spawn(
		constants.AGENT_TYPE_GRASS_HOPPER, agents.AgentCreationOptions{
			X: spawnPos.X + 10, Y: spawnPos.Y,
		},
	)
	World.Entities.Agents.Spawn(
		constants.AGENT_TYPE_ENEMY_SOLDIER, agents.AgentCreationOptions{
			X: spawnPos.X + 15, Y: spawnPos.Y,
		},
	)
	// playerAgentID = World.Entities.Agents.Spawn(
	// 	constants.AGENT_TYPE_PLAYER, agents.AgentCreationOptions{
	// 		X: spawnPos.X, Y: spawnPos.Y,
	// 	},
	// )
	//	player = game.findPlayer()
	player.InventoryID = item.NewDevInventory()
	components.Init(&World, Cam, player)

	sound.LoadSound("player_jump", "jump.wav")
	Console = console.New()

	//add oxygen emitter
	//	particle_emitter.EmitOxygen(playerAgentID, &World.Entities.Particles)
	render.NewColorShader()
}

func keyCallBack(w *glfw.Window, k glfw.Key, s int, a glfw.Action, mk glfw.ModifierKey) {
	if a == glfw.Press && k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
}

func main() {
	log.Printf("done")
	init1()
	// win := render.NewWindow(800, 600, true)
	// window := win.Window
	// window.SetKeyCallback(keyCallBack)
	// defer glfw.Terminate()
	// spriteloader.InitSpriteloader(&win)

	// test()
	// last := time.Now()

	// dt := float32(0)
	// lastFrame := float32(glfw.GetTime())

	// for !window.ShouldClose() {
	// 	// dt := time.Since(last).Seconds()
	// 	// last = time.Now()
	// 	// center := cx.Vec{X:100,Y:100}
	// 	// center.Y = 100
	// 	// character.Update(dt, 250, 250)
	// 	// draw(window, program)

	// 	currTime := float32(glfw.GetTime())
	// 	dt = currTime - lastFrame
	// 	lastFrame = currTime

	// 	game.ProcessInput()
	// 	log.Printf("%d", dt)
	// }
}

func getPng(dir string) []string {
	var ret []string
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".png") {
				return nil
			}
			ret = append(ret, path)
			return nil
		})

	if err != nil {
		log.Printf("error dir %s\n", dir)
	}

	return ret
}

func test() {

	for _, loc := range animation.LoadList("./animation") {
		character, err := c.LoadCharacter(loc)
		if err != nil {
			log.Println(loc.Name, err)
			continue
		}

		for _, skin := range character.Skeleton.Data.Skins {
			for _, att := range skin.Attachments {
				if _, ismesh := att.(*spine.MeshAttachment); ismesh {
					log.Println(loc.Name, "Unsupported")
					//	continue skip
				}
			}
		}

		characters = append(characters, character)
	}

	// var character *c.Character
	characterIndex := 0
	character = characters[characterIndex]

	log.Println("%v", character)

}
