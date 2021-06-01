package render

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	Height    int
	Width     int
	Resizable bool
	Window    *glfw.Window
	Program   uint32
	VAO       uint32
}

func NewWindow(width, height int, resizable bool) Window {
	glfwWindow := initGlfw(width, height, resizable)
	program := initOpenGL()

	//temporary, to set projection matrix

	window := Window{
		Width:     width,
		Height:    height,
		Resizable: resizable,
		Window:    glfwWindow,
		Program:   program,
		VAO:       makeVao(),
	}

	projectionMatrix := window.GetProjectionMatrix()
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projectionMatrix[0])

	return window
}

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

	//unbind
	gl.BindVertexArray(0)

	return vao
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

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	var vertexShaderSource string
	var fragmentShaderSource string

	vertexShaderSource = `
	#version 410
	layout (location=0) in vec3 position;
	layout (location=1) in vec2 texcoord;
	out vec2 tCoord;
	uniform mat4 projection;
	uniform mat4 world;
	uniform mat4 view;
	uniform vec2 texScale;
	uniform vec2 texOffset;
	void main() {
		gl_Position = projection  *  world * vec4(position, 1.0);
		tCoord = (texcoord+texOffset) * texScale;
	}
	` + "\x00"
	//gl_Position = vec4(position, 10.0, 1.0) * camera * projection;

	fragmentShaderSource = `
	#version 410
	in vec2 tCoord;
	out vec4 frag_colour;
	uniform sampler2D ourTexture;
	uniform vec4 color;
	void main() {
			frag_colour = texture(ourTexture, tCoord) * color;
	}
	` + "\x00"

	prog := CreateProgram(vertexShaderSource, fragmentShaderSource)

	gl.UseProgram(prog)
	gl.Uniform2f(
		gl.GetUniformLocation(prog, gl.Str("texScale\x00")),
		1.0, 1.0,
	)
	gl.Uniform4f(
		gl.GetUniformLocation(prog, gl.Str("color\x00")),
		1, 1, 1, 1,
	)

	// line opengl program
	vertexShaderSource = `
	#version 330 core
	layout (location = 0) in vec3 aPos;
	uniform mat4 uProjection;
	uniform mat4 uWorld;

	void main()
	{
	   gl_Position = uProjection * vec4(aPos, 1.0);
	}` + "\x00"

	fragmentShaderSource = `
	#version 330 core
	out vec4 FragColor;
	uniform vec3 uColor;

	void main()
	{
	   FragColor = vec4(uColor, 1.0f);
	}` + "\x00"

	lineProgram = CreateProgram(vertexShaderSource, fragmentShaderSource)

	return prog
}

func CreateProgram(vertexShaderSource string, fragmentShaderSource string) uint32 {
	vertexShader, err := CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	gl.DeleteShader(fragmentShader)
	gl.DeleteShader(vertexShader)

	return prog
}

func CompileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func (window *Window) GetProjectionMatrix() mgl32.Mat4 {
	return window.DefaultRenderContext().Projection
}
