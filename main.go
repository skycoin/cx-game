package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/starmap"

	//cv "github.com/skycoin/cx-game/cmd/spritetool"

	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
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
	sprite = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,

		1, -1, 0, 1, 1,
		-1, -1, 0, 0, 1,
		-1, 1, 0, 0, 0,
	}
)

var wz float32
var upPressed bool
var downPressed bool
var leftPressed bool
var rightPressed bool
var spacePressed bool

var isFreeCam = false

var cat *models.Cat
var fps *models.Fps

var Cam *camera.Camera
var tex uint32

func makeVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(sprite), gl.Ptr(sprite), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)

	return vao
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
			wz += 0.5
		}
		if k == glfw.KeyZ {
			wz -= 0.5
		}
		if k == glfw.KeyF2 {
			isFreeCam = !isFreeCam
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

var win render.Window

func main() {

	/*
		var SS cv.SpriteSet
		SS.LoadFile("./assets/sprite.png", 250, false)
		SS.ProcessContours()
		SS.DrawSprite()
	*/

	win = render.NewWindow(height, width, true)
	spriteloader.InitSpriteloader(&win)
	cat = models.NewCat()
	fps = models.NewFps(false)

	wz = -10
	CurrentPlanet = world.NewDevPlanet()
	window := win.Window
	Cam = camera.NewCamera(&win)
	spawnX := int(20)
	Cam.X = float32(spawnX)
	Cam.Y = 5
	cat.Pos.X = float32(spawnX)
	cat.Pos.Y = float32(CurrentPlanet.GetHeight(spawnX) + 6)

	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	VAO := makeVao()
	program := win.Program
	gl.GenTextures(1, &tex)
	lastTime := models.GetTimeStamp()

	starmap.Init(&win)
	starmap.Generate(256, 0.04, 8)

	for !window.ShouldClose() {
		currTime := models.GetTimeStamp()
		elapsed := currTime - lastTime
		Tick(elapsed)
		redraw(window, program, VAO)
		fps.Tick()
		lastTime = currTime
	}
}

func boolToFloat(x bool) float32 {
	if x {
		return 1
	} else {
		return 0
	}
}

func Tick(elapsed int) {
	dt := float32(elapsed) / 1000

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

	starmap.Draw()
	CurrentPlanet.Draw(Cam)
	cat.Draw(Cam)

	// tile - air line (green)
	collidingTileLines := CurrentPlanet.GetCollidingTilesLinesRelative(
		int(cat.Pos.X), int(cat.Pos.Y), Cam)
	if len(collidingTileLines) > 2 {
		win.DrawLines(collidingTileLines, []float32{0.0, 1.0, 0.0})
	}

	// body bounding box (blue)
	win.DrawLines(cat.GetBBoxLinesRelative(Cam), []float32{0.0, 0.0, 1.0})

	// colliding line from body (red)
	collidingLines := cat.GetCollidingLinesRelative(Cam)
	if len(collidingLines) > 2 {
		win.DrawLines(collidingLines, []float32{1.0, 0.0, 0.0})
	}

	glfw.PollEvents()
	window.SwapBuffers()
}
