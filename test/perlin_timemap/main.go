package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
)

var (
	vertices []float32
	myPerlin perlin.Perlin2D
)

func main() {
	win := render.NewWindow(800, 600, false)
	glfwWindow := win.Window

	shader := utility.NewShader("./test/perlin_timemap/shaders/vertex.glsl", "./test/perlin_timemap/shaders/fragment.glsl")
	shader.Use()
	InitPerlin2D()

	ortho := mgl32.Ortho2D(0, 1000, -1, 1)
	shader.SetMat4("ortho", &ortho)

	var dt, lastFrame float32
	lastFrame = float32(glfw.GetTime())
	// var x float32
	for !glfwWindow.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// gl.BindVertexArray(vao)
		// gl.DrawArrays(gl.TRIANGLES, 0, 6)
		currentFrame := float32(glfw.GetTime())
		dt = currentFrame - lastFrame
		lastFrame = currentFrame

		// x = x + 200*dt
		// world := mgl32.Translate3D(
		// 	x,
		// 	0,
		// 	0,
		// )
		// shader.SetMat4("world", &world)
		z := float32(glfw.GetTime())

		GenData(GenPerlin(1, z))

		ClearValues()
		DrawPerlin()

		glfw.PollEvents()
		glfwWindow.SwapBuffers()
	}
	dt = dt + 1
}

var localPoints []float32
var counter float32 = 1

func GenData(noise float32) {
	if len(localPoints) < 100 {
		localPoints = append(localPoints, counter, noise)
		counter++
		return
	}

	GenVAO(localPoints)
	localPoints = nil
}

func GenVAO(points []float32) {
	var vao, vbo uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(points)*4, gl.Ptr(points), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	vaos = append(vaos, vao)
	vbos = append(vbos, vbo)
}

var vaos []uint32
var vbos []uint32

func ClearValues() {
	if len(vaos) < 16 {
		return
	}

	for _, vbo := range vbos {
		gl.DeleteBuffers(1, &vbo)
	}
	for _, vao := range vaos {
		gl.DeleteVertexArrays(1, &vao)
	}

	vaos = nil
	vbos = nil
}
func DrawPerlin() {
	for _, vao := range vaos {
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.LINE_STRIP, 0, 100)
	}
}
func InitPerlin2D() {
	myPerlin = perlin.NewPerlin2D(
		3,
		512,
		4,
		256,
	)
}
func InitPerlin3D() {

}

func GenPerlin(x, y float32) float32 {
	result := myPerlin.Noise(
		x,
		y,
		0.2,
		3,
		8,
	)

	return result
}
