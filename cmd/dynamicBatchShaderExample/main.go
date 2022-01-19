package main

import (
	"fmt"
	_ "image/png"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	indexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBuffer"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/Texture"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/UI_Injector"
	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBuffer"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/renderer"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/shader"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexArray"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexbufferLayout"
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
	width  = 960
	height = 540
)

var (
	positions = []float32{
		-50.0, -50.0, 0.0, 0.0, //1
		50.0, -50.0, 1.0, 0.0, //2
		50.0, 50.0, 1.0, 1.0, // 3
		-50.0, 50.0, 0.0, 1.0, //4

	}

	indices = []uint32{
		0, 1, 2,
		2, 3, 0,
	}
)

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 0)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "batch", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() {
	// // Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
}

var ib *indexbuffer.IndexBuffer
var vb *vertexbuffer.VertexBuffer
var va *vertexArray.VertexArray
var vbl *vertexbufferLayout.VertexbufferLayout

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	initOpenGL()
	var UI *UI_Injector.UI_Injector

	UI = UI_Injector.SetUpUI()
	go UI.ListenForChanges()
	var objectAdjustment = UI
	var proj mgl32.Mat4 = mgl32.Ortho(0.0, 960.0, 0.0, 540.0, -1.0, 1.0)
	var view mgl32.Mat4 = mgl32.Translate3D(0.0, 0.0, 0.0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	shader := shader.SetupShader("./../../assets/shader/spine/basic.shader")
	shader.Bind()
	shader.SetUniForm4f("u_Color", 0.8, 0.3, 0.8, 1.0)

	shader.UnBind()

	//setup vertex array
	va = vertexArray.SetUpVertxArray()
	// setup and run vertex buffer
	vb = vertexbuffer.RunVertexBuffer(positions, len(positions)*4*4)
	//setup vertex layout
	vbl = &vertexbufferLayout.VertexbufferLayout{}
	//add vertex buffer to vertex bufferlayout
	vbl.Push(gl.FLOAT, 2)
	vbl.Push(gl.FLOAT, 2)
	va.AddBuffer(vb, vbl)
	// setup and run index buffer
	ib = indexbuffer.RunIndexBuffer(indices, 6)

	tex := Texture.SetUpTexture("./cat.png")
	tex.Bind(0)
	shader.SetUniForm1i("u_Texture", 0)

	va.Unbind()
	vb.Unbind()
	ib.Unbind()

	var translationA = mgl32.Translate3D(objectAdjustment.Object[0].X, objectAdjustment.Object[0].Y, objectAdjustment.Object[0].Z)

	var translationB = mgl32.Translate3D(400+objectAdjustment.Object[1].X, 200+objectAdjustment.Object[1].Y, objectAdjustment.Object[1].Z)

	render := renderer.SetupRender()
	var r float32 = 0.0
	var increment float32 = 0.5

	for !window.ShouldClose() {
		render.Clear()

		//fmt.Println(translation)
		// shader.SetUniForm4f("u_Color", r, 0.3, 0.8, 1.0)

		{
			var model mgl32.Mat4 = translationA.Add(mgl32.Translate3D(objectAdjustment.Object[0].X*10, objectAdjustment.Object[0].Y*10, objectAdjustment.Object[0].Z*10))
			mvp := proj.Mul4(view).Mul4(model)
			shader.Bind()
			shader.SetUniFormMat4f("u_MVP", mvp)
			render.Draw(va, ib, shader)
		}

		{
			var model mgl32.Mat4 = translationB.Add(mgl32.Translate3D(objectAdjustment.Object[1].X*10, objectAdjustment.Object[1].Y*10, objectAdjustment.Object[1].Z*10))
			mvp := proj.Mul4(view).Mul4(model)
			shader.Bind()
			shader.SetUniFormMat4f("u_MVP", mvp)
			render.Draw(va, ib, shader)
		}

		if r > 1.9 {
			increment = -0.05
		} else if r < 0.0 {
			increment = 0.05
		}
		r = r + increment
		glfw.PollEvents()
		window.SwapBuffers()
	}

}
