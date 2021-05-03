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
    void main() {
		frag_colour = texture(ourTexture, tCoord);
    }
` + "\x00"

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
	gl.UseProgram(prog)
	gl.Uniform2f(
		gl.GetUniformLocation(prog, gl.Str("texScale\x00")),
		1.0, 1.0,
	)
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

	vertexShader, err := CompileShader(string(vertexShaderSource)+"\x00", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := CompileShader(string(fragmentShaderSource)+"\x00", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	gl.UseProgram(prog)
	return prog
}
