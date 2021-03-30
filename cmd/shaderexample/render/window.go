package render

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"strings"
)

type Window struct {
	Height    int
	Width     int
	Resizable bool
	Window    *glfw.Window
	Program   uint32
}

func NewWindow(height, width int, resizable bool) (window Window, closeFn func()) {
	glWindow, glWindowCloseFn := initGlfw(height, width, resizable)
	program, err := newProgram()
	if err != nil {
		panic(err)
	}

	window = Window{
		Height:    height,
		Width:     width,
		Resizable: resizable,
		Window:    glWindow,
		Program:   program,
	}

	closeFn = func() {
		glWindowCloseFn()
		// add another close instructions
	}

	return window, closeFn
}

// initGlfw initializes glfw and returns a Window to use.
func initGlfw(height, width int, resizable bool) (*glfw.Window, func()) {
	if err := glfw.Init(); err != nil {
		panic("failed to initialize glfw")
	}

	var res int
	if resizable {
		res = glfw.True
	} else {
		res = glfw.False
	}

	glfw.WindowHint(glfw.Resizable, res)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)
	window, err := glfw.CreateWindow(width, height, "CX Game", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	cancelFn := func() {
		glfw.Terminate()
	}
	return window, cancelFn
}

func newProgram() (uint32, error) {
	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

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

		return 0, errors.New(fmt.Sprintf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	cStrSource := gl.Str(source)
	gl.ShaderSource(shader, 1, &cStrSource, nil)
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
