package shader

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	Filepath               string
	M_renderID             uint32
	M_UniformLocationCache map[string]int32
}

type ShaderProgramSource struct {
	VertexSource   string
	FragmentSource string
}

func SetupShader(filepath string) *Shader {
	newShader := &Shader{}
	newShader.Filepath = filepath
	newShader.M_renderID = 0

	s := newShader.ParseShader(filepath)
	program, err := newShader.CreateShader(s.VertexSource, s.FragmentSource)
	if err != nil {
		panic(err)
	}

	newShader.M_renderID = program
	newShader.M_UniformLocationCache = make(map[string]int32)
	return newShader
}

func (s *Shader) SetUniForm4f(name string, v0, v1, v2, v3 float32) {
	gl.Uniform4f(s.getUniformLocation(name), v0, v1, v2, v3)
}
func (s *Shader) SetUniFormMat4f(name string, matrix mgl32.Mat4) {
	gl.UniformMatrix4fv(s.getUniformLocation(name), 1, false, &matrix[0])
}
func (s *Shader) SetUniForm1f(name string, value float32) {
	gl.Uniform1f(s.getUniformLocation(name), value)
}
func (s *Shader) SetUniForm1i(name string, value int32) {
	gl.Uniform1i(s.getUniformLocation(name), value)
}
func (s *Shader) SetUniform1iv(name string, count int32, value *int32) {
	gl.Uniform1iv(s.getUniformLocation(name), count, value)
}

func (s *Shader) getUniformLocation(name string) int32 {

	_, ok := s.M_UniformLocationCache[name]
	if ok {
		//	panic("already got uniform")
		return s.M_UniformLocationCache[name]
	}
	location := gl.GetUniformLocation(s.M_renderID, gl.Str(name+"\x00"))
	if location == -1 {
		fmt.Println("uniform: " + name + " does not exist!! ")
	}

	s.M_UniformLocationCache[name] = location
	//fmt.Println("got first uniform name")

	return location
}

func (s *Shader) ParseShader(filepath string) *ShaderProgramSource {

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

func (s *Shader) CreateShader(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := s.compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fmt.Println("done vertex shader")

	fragmentShader, err := s.compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
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

func (s *Shader) compileShader(source string, shaderType uint32) (uint32, error) {
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

func (s *Shader) DeleteShader() {
	gl.DeleteProgram(s.M_renderID)
}

func (s *Shader) Bind() {
	gl.UseProgram(s.M_renderID)
}

func (s *Shader) UnBind() {
	gl.UseProgram(0)
}
