package render

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
		VAO:       0,
	}
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
    void main() {
		
        gl_Position = projection * world * vec4(position, 1.0);
		tCoord = texcoord;
    }
` + "\x00"
	//gl_Position = vec4(position, 10.0, 1.0) * camera * projection;

	fragmentShaderSource = `
    #version 410
	in vec2 tCoord;
    out vec4 frag_colour1;
	uniform sampler2D textures[8];
    void main() {
		 frag_colour1 = texture(textures[1], tCoord);
    }
` + "\x00"

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
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
