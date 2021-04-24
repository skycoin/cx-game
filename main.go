package main

import (
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"

	//cv "github.com/skycoin/cx-game/cmd/spritetool"

	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/render"
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

	gravity = 0.01
)

var (
	sprite = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,
		-1, -1, 0, 0, 1,
	}
	indices = []int32{
		1, 0, 2,
		1, 2, 3,
	}
)

var wx, wy, wz float32
var upPressed bool
var downPressed bool
var leftPressed bool
var rightPressed bool
var spacePressed bool

var cat *models.Cat
var fps *models.Fps

var Cam camera.Camera
var tex uint32

func makeVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var ebo uint32
	gl.GenBuffers(1, &ebo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(sprite), gl.Ptr(sprite), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
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
	wy = 0
	wz = -10
	win := render.NewWindow(height, width, true)
	window := win.Window
	window.SetKeyCallback(keyCallBack)
	defer glfw.Terminate()
	VAO := makeVao()
	program := win.Program
	gl.GenTextures(1, &tex)
	for !window.ShouldClose() {
		Tick()
		redraw(window, program, VAO)
		fps.Tick()
	}
}

func Tick() {
	if wy > -3 {
		cat.YVelocity -= gravity
	} else {
		cat.YVelocity = 0

		if wy > -3 {
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

		wx += cat.XVelocity
		wy += cat.YVelocity

		spacePressed = false
	}
}

func redraw(window *glfw.Window, program uint32, VAO uint32) {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// cat := models.NewCat()
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(cat.Size.X), int32(cat.Size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(cat.RGBA.Pix))
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture\x00")), 0)
	worldTranslate := mgl32.Translate3D(wx, wy, wz)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("world\x00")), 1, false, &worldTranslate[0])
	projectTransform := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projectTransform[0])
	gl.BindVertexArray(VAO)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	glfw.PollEvents()
	window.SwapBuffers()
}
