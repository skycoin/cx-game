// package main_test

// import (
// 	"bufio"
// 	"fmt"
// 	_ "image/png"
// 	"log"
// 	"os"
// 	"runtime"
// 	"strings"

// 	"github.com/go-gl/gl/v4.1-core/gl"
// 	"github.com/go-gl/glfw/v3.3/glfw"
// 	indexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBuffer"
// 	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBuffer"
// 	"github.com/skycoin/cx-game/world"
// )

// func init() {
// 	// This is needed to arrange that main() runs on main thread.
// 	// See documentation for functions that are only allowed to be called from the main thread.
// 	runtime.LockOSThread()
// }

// var (
// 	DrawCollisionBoxes = false
// 	FPS                int
// )

// var CurrentPlanet *world.Planet

// const (
// 	width  = 800
// 	height = 480
// )

// var (
// 	positions = []float32{
// 		-0.5, -0.5, //1
// 		0.5, -0.5, //2
// 		0.5, 0.5, // 3
// 		-0.5, 0.5, //4

// 	}

// 	indices = []uint32{
// 		0, 1, 2,
// 		2, 3, 0,
// 	}
// )

// func initGlfw() *glfw.Window {
// 	if err := glfw.Init(); err != nil {
// 		panic(err)
// 	}
// 	glfw.WindowHint(glfw.Resizable, glfw.False)
// 	glfw.WindowHint(glfw.ContextVersionMajor, 4)
// 	glfw.WindowHint(glfw.ContextVersionMinor, 0)
// 	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
// 	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

// 	window, err := glfw.CreateWindow(width, height, "batch", nil, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	window.MakeContextCurrent()
// 	glfw.SwapInterval(1)

// 	return window
// }

// // initOpenGL initializes OpenGL and returns an intiialized program.
// func initOpenGL() {
// 	// s := ParseShader("./../../assets/shader/spine/basic.shader")
// 	// fmt.Println("shader code: ", s.FragmentSource)
// 	// fmt.Println("shader code: ", s.VertexSource)
// 	// Initialize Glow
// 	if err := gl.Init(); err != nil {
// 		panic(err)
// 	}

// 	version := gl.GoStr(gl.GetString(gl.VERSION))
// 	fmt.Println("OpenGL version", version)

// 	// Configure the vertex and fragment shaders
// 	// program, err := newProgram(s.VertexSource, s.FragmentSource)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// return program
// }

// var ib *indexbuffer.IndexBuffer
// var vb *vertexbuffer.VertexBuffer

// // makeVao initializes and returns a vertex array from the points provided.
// func makeVao(points []float32) uint32 {

// 	var vao uint32
// 	gl.GenVertexArrays(1, &vao)
// 	gl.BindVertexArray(vao)

// 	var vbo uint32
// 	// vb = vertexbuffer.RunVertexBuffer(points, 2*4*len(points))
// 	gl.GenBuffers(1, &vbo)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
// 	gl.BufferData(gl.ARRAY_BUFFER, 2*6*len(points), gl.Ptr(points), gl.STATIC_DRAW)

// 	gl.EnableVertexAttribArray(0)
// 	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*int32(len(points)), gl.PtrOffset(0))

// 	/// index buffer
// 	var ibo uint32

// 	//ib = indexbuffer.RunIndexBuffer(indices, 6)
// 	gl.GenBuffers(1, &ibo)
// 	//	fmt.Println("this is the ID Buffer: ", IB.M_renderID)
// 	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
// 	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 6*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

// 	return ibo
// }

// func main() {
// 	runtime.LockOSThread()

// 	window := initGlfw()
// 	defer glfw.Terminate()
// 	initOpenGL()
// 	s := ParseShader("./../../assets/shader/spine/basic.shader")
// 	program, err := newProgram(s.VertexSource, s.FragmentSource)
// 	if err != nil {
// 		panic(err)
// 	}

// 	gl.UseProgram(program)

// 	//	return
// 	//renderer.GLClearError()
// 	//vao := makeVao(positions)
// 	var vao uint32
// 	gl.GenVertexArrays(1, &vao)
// 	gl.BindVertexArray(vao)

// 	//  unsafe.Sizeof(float32(0.0))

// 	var buffer uint32
// 	// vb = vertexbuffer.RunVertexBuffer(points, 2*4*len(points))
// 	gl.GenBuffers(1, &buffer)
// 	gl.BindBuffer(gl.ARRAY_BUFFER, buffer)

// 	// 4 = sizeof(float)  2=number of triangles  len(position) = total number of elements inside position
// 	gl.BufferData(gl.ARRAY_BUFFER, len(positions)*2*4, gl.Ptr(positions), gl.STATIC_DRAW)
// 	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
// 	gl.EnableVertexAttribArray(0)

// 	/// index buffer
// 	var ibo uint32

// 	//ib = indexbuffer.RunIndexBuffer(indices, 6)
// 	gl.GenBuffers(1, &ibo)
// 	// //	fmt.Println("this is the ID Buffer: ", IB.M_renderID)
// 	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ibo)
// 	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

// 	// fmt.Println(vao)
// 	// var buffer uint32
// 	// gl.GenBuffers(1, &buffer)
// 	// gl.BindBuffer(gl.ARRAY_BUFFER, buffer)
// 	// gl.BufferData(gl.ARRAY_BUFFER, len(positions)*6, gl.Ptr(positions), gl.STATIC_DRAW)

// 	// gl.EnableVertexAttribArray(0)
// 	// gl.VertexAttribPointer(0, 2, gl.FLOAT, false, int32(len(positions))*2, gl.PtrOffset(0))

// 	//gl.BindBuffer(gl.ARRAY_BUFFER, 0)

// 	location := gl.GetUniformLocation(program, gl.Str("u_Color\x00"))

// 	var r float32 = 0.0
// 	var increment float32 = 0.5
// 	fmt.Println("uniform location: ", location)

// 	// fmt.Println("number check: ",)

// 	for !window.ShouldClose() {
// 		gl.Clear(gl.COLOR_BUFFER_BIT)

// 		//	gl.UseProgram(program)

// 		//gl.BindVertexArray(ibo)
// 		//ib.Bind()
// 		gl.Uniform4f(location, r, 0.3, 0.8, 1.0)
// 		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
// 		//	renderer.GLCheckError()
// 		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(positions)))

// 		if r > 1.9 {
// 			increment = -0.05
// 		} else if r < 0.0 {
// 			increment = 0.05
// 		}
// 		r = r + increment
// 		glfw.PollEvents()
// 		window.SwapBuffers()
// 	}

// 	gl.DeleteProgram(program)
// }

// type ShaderProgramSource struct {
// 	VertexSource   string
// 	FragmentSource string
// }

// func ParseShader(filepath string) *ShaderProgramSource {

// 	type ShaderType int64

// 	const (
// 		NONE     ShaderType = -1
// 		VERTEX   ShaderType = 0
// 		FRAGMENT ShaderType = 1
// 	)

// 	var currentShaderType ShaderType = NONE

// 	var shaderStream [2]string

// 	fmt.Println("currnet shader type: ", currentShaderType)

// 	f, err := os.Open(filepath) // Error handling elided for brevity.
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)

// 	for scanner.Scan() {

// 		if strings.Contains(scanner.Text(), "#shader") {

// 			if strings.Contains(scanner.Text(), "vertex") {
// 				currentShaderType = VERTEX
// 			} else if strings.Contains(scanner.Text(), "fragment") {
// 				currentShaderType = FRAGMENT
// 			}
// 		} else {
// 			shaderStream[int64(currentShaderType)] += (scanner.Text() + "\n")
// 		}
// 	}
// 	if err := scanner.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	shaderStream[0] += "\x00"
// 	shaderStream[1] += "\x00"

// 	ss := ShaderProgramSource{shaderStream[0], shaderStream[1]}
// 	return &ss
// }
// func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
// 	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
// 	if err != nil {
// 		return 0, err
// 	}

// 	fmt.Println("done vertex shader")

// 	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
// 	if err != nil {

// 		return 0, err
// 	}

// 	fmt.Println("done fragment shader")

// 	program := gl.CreateProgram()

// 	gl.AttachShader(program, vertexShader)
// 	gl.AttachShader(program, fragmentShader)
// 	gl.LinkProgram(program)

// 	var status int32
// 	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
// 	if status == gl.FALSE {
// 		var logLength int32
// 		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

// 		log := strings.Repeat("\x00", int(logLength+1))
// 		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

// 		return 0, fmt.Errorf("failed to link program: %v", log)
// 	}

// 	gl.DeleteShader(vertexShader)
// 	gl.DeleteShader(fragmentShader)

// 	return program, nil
// }

// func compileShader(source string, shaderType uint32) (uint32, error) {
// 	shader := gl.CreateShader(shaderType)

// 	csources, _ := gl.Strs(source)
// 	gl.ShaderSource(shader, 1, csources, nil)
// 	//	free()
// 	gl.CompileShader(shader)

// 	var status int32
// 	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
// 	if status == gl.FALSE {
// 		var logLength int32
// 		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

// 		log := strings.Repeat("\x00", int(logLength+1))
// 		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
// 		gl.DeleteShader(shader)

// 		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
// 	}

// 	return shader, nil
// }
