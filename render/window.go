package render

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/engine/input"
)

type Window struct {
	// note that the virtual dimensions have no prefix.
	// i.e. "Window.Width" is really the virtual width
	// it is assumed that the vast majority of the time,
	// the virtual dimensions are what the programmer wants.
	Width, Height                 int
	PhysicalWidth, PhysicalHeight int
	Resizable                     bool
	Window                        *glfw.Window
	context                       Context

	// used for mouse events
	PhysicalToViewportTransform mgl32.Mat4
}

func (win *Window) sizeCallback(
	window *glfw.Window, physicalWidth, physicalHeight int,
) {
	// "physical" dimensions describe actual window size
	// "virtual" dimensions describe scaling of both world and UI
	// physical determines resolution.
	// virtual determines how big things are.
	virtualWidth := float32(win.Width)
	virtualHeight := float32(win.Height)
	windowDimensions := fitCentered(
		mgl32.Vec2{virtualWidth, virtualHeight},
		mgl32.Vec2{float32(physicalWidth), float32(physicalHeight)},
	)
	windowDimensions.Viewport.Use()

	win.PhysicalToViewportTransform = windowDimensions.Transform()

	input.SetPhysicalToViewportTransform(win.PhysicalToViewportTransform)

	Projection = mgl32.Ortho(
		-virtualWidth/2/PixelsPerTile, virtualWidth/2/PixelsPerTile,
		-virtualHeight/2/PixelsPerTile, virtualHeight/2/PixelsPerTile,
		-1000, 1000,
	)
}

func (win *Window) SetInitialWindowDimensions() {
	virtualWidth := float32(win.Width)
	virtualHeight := float32(win.Height)
	windowDimensions := fitCentered(
		mgl32.Vec2{virtualWidth, virtualHeight},
		mgl32.Vec2{virtualWidth, virtualHeight},
	)
	windowDimensions.Viewport.Use()

	win.PhysicalToViewportTransform = windowDimensions.Transform()
	input.SetPhysicalToViewportTransform(win.PhysicalToViewportTransform)

}

func NewWindow(width, height int, resizable bool) Window {
	glfwWindow := initGlfw(width, height, resizable)
	initOpenGL()

	InitQuadVao()

	//temporary, to set projection matrix

	window := Window{
		Width:     width,
		Height:    height,
		Resizable: resizable,
		Window:    glfwWindow,
	}
	Projection = mgl32.Ortho(
		-float32(window.Width)/2/PixelsPerTile, float32(window.Width)/2/PixelsPerTile,
		-float32(window.Height)/2/PixelsPerTile, float32(window.Height)/2/PixelsPerTile,
		-1000, 1000,
	)
	window.context = window.DefaultRenderContext()

	return window
}

func (w *Window) SetCallbacks() {
	w.Window.SetSizeCallback(w.sizeCallback)
	w.SetInitialWindowDimensions()
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(width, height int, resizable bool) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	var res int
	if resizable {
		res = glfw.True
	} else {
		res = glfw.False
	}

	glfw.WindowHint(glfw.Resizable, res)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "CX Game", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// https://github.com/glfw/glfw/issues/1334
func FixRenderCOCOA(window *glfw.Window) {
	windowMoved := false
	if !windowMoved {
		x, y := window.GetPos()
		window.SetPos(x+1, y)
		windowMoved = true
	}
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	InitDrawLines()
	lineProgram = CompileProgram(
		"./assets/shader/line.vert", "./assets/shader/line.frag")
}

var Projection mgl32.Mat4

func (window *Window) GetProjectionMatrix() mgl32.Mat4 {
	// return window.DefaultRenderContext().Projection
	return window.context.Projection
}

func (window *Window) SetProjectionMatrix(projection mgl32.Mat4) {
	window.context.Projection = projection
	Projection = projection
}
