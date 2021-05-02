package main

import (
    "log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/starmap"

	//cv "github.com/skycoin/cx-game/cmd/spritetool"

	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/ui"
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

	gravity = 0.01
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

var wx, wy, wz float32
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
    if a != glfw.Press {return}
    screenX := float32(2*mouseX/float64(win.Width)-1)
    screenY := 1-float32(2*mouseY/float64(win.Height))
    projection := win.GetProjectionMatrix()

    didSelectPaleteTile := tilePaleteSelector.
        TrySelectTile(screenX,screenY,projection)

    log.Printf("selected paleete tile? %v", didSelectPaleteTile)
    if !didSelectPaleteTile {
        CurrentPlanet.TryPlaceTile(
            screenX,screenY,
            projection,
            tilePaleteSelector.GetSelectedTile(),
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
var tilePaleteSelector ui.TilePaleteSelector

var cat *models.Cat
var fps *models.Fps

var Cam *camera.Camera
var win render.Window
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
			//wy += 0.5
			upPressed = true
		}
		if k == glfw.KeyS {
			// wy -= 0.5
			downPressed = true
		}
		if k == glfw.KeyA {
			// wx -= 0.5
			leftPressed = true
		}
		if k == glfw.KeyD {
			// wx += 0.5
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
		if k == glfw.KeyF3 {
			isTileSelectorVisible = !isTileSelectorVisible
		}
	} else if a == glfw.Release {
		if k == glfw.KeyW {
			//wy += 0.5
			upPressed = false
		}
		if k == glfw.KeyS {
			// wy -= 0.5
			downPressed = false
		}
		if k == glfw.KeyA {
			// wx -= 0.5
			leftPressed = false
		}
		if k == glfw.KeyD {
			// wx += 0.5
			rightPressed = false
		}
	}
}

func main() {

	/*
		var SS cv.SpriteSet
		SS.LoadFile("./assets/sprite.png", 250, false)
		SS.ProcessContours()
		SS.DrawSprite()
	*/

	cat = models.NewCat()
	fps = models.NewFps(false)

	wx = 0
	wy = 10
	wz = -10
	win = render.NewWindow(height, width, true)
	spriteloader.InitSpriteloader(&win)
	CurrentPlanet = world.NewDevPlanet()
    worldTiles := CurrentPlanet.GetAllTilesUnique()
    log.Printf("Found [%v] unique tiles in the world",len(worldTiles))
    tilePaleteSelector = ui.
        MakeTilePaleteSelector(worldTiles)
	window := win.Window
	Cam = camera.NewCamera(&win)
	wx = 20
	Cam.X = 20
	Cam.Y = 5

	window.SetKeyCallback(keyCallBack)
    window.SetCursorPosCallback(cursorPosCallback)
    window.SetMouseButtonCallback(mouseButtonCallback)
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

func boolToInt(x bool) int {
	if x {
		return 1
	} else {
		return 0
	}
}

func Tick(elapsed int) {
	if wy > 6.5 {
		cat.YVelocity -= gravity
	} else {
		cat.YVelocity = 0

		if spacePressed {
			cat.YVelocity = 0.2
		}
	}

	if !rightPressed || !leftPressed {
		cat.XVelocity = 0
	}

	if rightPressed {
		cat.XVelocity = 0.05
	}

	if leftPressed {
		cat.XVelocity = -0.05
	}

	if isFreeCam {
		Cam.MoveCam(
			float32(boolToInt(rightPressed)-boolToInt(leftPressed)),
			float32(boolToInt(upPressed)-boolToInt(downPressed)),
			0,
			float32(elapsed)/1000,
		)
	} else {
		wx += cat.XVelocity
		wy += cat.YVelocity
	}

	spacePressed = false
}

func redraw(window *glfw.Window, program uint32, VAO uint32) {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	starmap.Draw()

	gl.UseProgram(program)

	// cat := models.NewCat()
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	//gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)
	gl.Disable(gl.DEPTH_TEST)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(cat.Size.X), int32(cat.Size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(cat.RGBA.Pix))

	gl.BindVertexArray(VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

	CurrentPlanet.Draw(Cam)

    if isTileSelectorVisible {
        tilePaleteSelector.Draw()
    }

	glfw.PollEvents()
	window.SwapBuffers()
}
