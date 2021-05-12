package render

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"path/filepath"
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

func NewWindow(height, width int, resizable bool) Window {
	window := initGlfw(height, width, resizable)
	program := initOpenGL()
	return Window{
		Height:    height,
		Width:     width,
		Resizable: resizable,
		Window:    window,
		Program:   program,
		VAO:       makeVao(),
	}
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

	return vao
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(height, width int, resizable bool) *glfw.Window {
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
	uniform vec2 texScale;
	uniform vec2 texOffset;
	void main() {
		gl_Position = projection * world * vec4(position, 1.0);
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
		1,1,1,1,
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
	aspect := float32(window.Width) / float32(window.Height)
	return mgl32.Perspective(
		mgl32.DegToRad(45), aspect, 0.1, 100.0,
	)
}

//opengl with custom shaders
func InitOpenGLCustom(shaderDir string) uint32 {

	var vertexShaderSource, fragmentShaderSource []byte

	filepath.WalkDir(shaderDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Panic(err)
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "vertex.glsl" {
			vertexShaderSource, err = ioutil.ReadFile(path)
			if err != nil {
				log.Panic(err)
			}
		} else if d.Name() == "fragment.glsl" {
			fragmentShaderSource, err = ioutil.ReadFile(path)
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if len(vertexShaderSource) == 0 || len(fragmentShaderSource) == 0 {
		log.Panic("NO shaders found")
	}
	if err := gl.Init(); err != nil {
		panic(err)
	}

	prog := CreateProgram(string(vertexShaderSource), string(fragmentShaderSource))
	gl.UseProgram(prog)

	return prog
}
