package main

import (
	"bufio"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	indexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBuffer"
	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBuffer"
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
	positions = []float32{
		-0.5, -0.5, //1
		0.5, -0.5, //2
		0.5, 0.5, // 3
		-0.5, 0.5, //4

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
func initOpenGL() uint32 {
	s := ParseShader("./../../assets/shader/spine/basic.shader")
	fmt.Println("shader code: ", s.FragmentSource)
	fmt.Println("shader code: ", s.VertexSource)
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(s.VertexSource, s.FragmentSource)
	if err != nil {
		panic(err)
	}

	return program
}

var ib *indexbuffer.IndexBuffer
var vb *vertexbuffer.VertexBuffer

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) {

	// var vao uint32
	// gl.GenVertexArrays(1, &vao)
	// gl.BindVertexArray(vao)

	// //var vbo uint32
	// vb = vertexbuffer.RunVertexBuffer(points, len(points)*2*4)

	// gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	// gl.EnableVertexAttribArray(0)
	// /// index buffer
	// //var ibo uint32

	// ib = indexbuffer.RunIndexBuffer(indices, 6)

	// //return vao
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()
	program := initOpenGL()
	gl.UseProgram(program)
	location := gl.GetUniformLocation(program, gl.Str("u_Color\x00"))

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	//var vbo uint32
	vb = vertexbuffer.RunVertexBuffer(positions, len(positions)*2*4)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	/// index buffer
	//var ibo uint32

	ib = indexbuffer.RunIndexBuffer(indices, 6)

	var r float32 = 0.0
	var increment float32 = 0.5

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)
		gl.Uniform4f(location, r, 0.3, 0.8, 1.0)

		//	gl.BindVertexArray(vao)
		//ib.Bind()

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
		//	renderer.GLCheckError()
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(positions)/3))

		if r > 1.9 {
			increment = -0.05
		} else if r < 0.0 {
			increment = 0.05
		}
		r = r + increment
		glfw.PollEvents()
		window.SwapBuffers()
	}

	gl.DeleteProgram(program)
}

type ShaderProgramSource struct {
	VertexSource   string
	FragmentSource string
}

func ParseShader(filepath string) *ShaderProgramSource {

	type ShaderType int64

	const (
		NONE     ShaderType = -1
		VERTEX   ShaderType = 0
		FRAGMENT ShaderType = 1
	)

	var currentShaderType ShaderType = NONE

	var shaderStream [2]string

	fmt.Println("currnet shader type: ", currentShaderType)

	f, err := os.Open(filepath) // Error handling elided for brevity.
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {

		if strings.Contains(scanner.Text(), "#shader") {

			if strings.Contains(scanner.Text(), "vertex") {
				currentShaderType = VERTEX
			} else if strings.Contains(scanner.Text(), "fragment") {
				currentShaderType = FRAGMENT
			}
		} else {
			shaderStream[int64(currentShaderType)] += (scanner.Text() + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	shaderStream[0] += "\x00"
	shaderStream[1] += "\x00"

	ss := ShaderProgramSource{shaderStream[0], shaderStream[1]}
	return &ss
}
func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fmt.Println("done vertex shader")

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {

		return 0, err
	}

	fmt.Println("done fragment shader")

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, _ := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	//	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		gl.DeleteShader(shader)

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
