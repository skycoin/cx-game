package render

// Helper shader program class from learnopengl.com

import (
	"io/ioutil"
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const header string = "#version 410\n"

type Program uint32
func (p Program) gl() uint32 { return uint32(p) } // downcast

type Shader uint32
func (s Shader) gl() uint32 { return uint32(s) } // downcast

type ShaderConfig struct {
	vertexPath string
	fragmentPath string
	definitions map[string]string
}

func NewShaderConfig(vertexPath, fragmentPath string) ShaderConfig {
	return ShaderConfig {
		vertexPath: vertexPath, fragmentPath: fragmentPath,
		definitions: make(map[string]string),
	}
}

func CompileProgram(vertexPath, fragmentPath string) Program {
	config := NewShaderConfig(vertexPath, fragmentPath)
	return config.Compile()
}

func (config *ShaderConfig) Define(name,value string) {
	config.definitions[name] = value
}

func (s Shader) Delete() {
	gl.DeleteShader(s.gl())
}

func (config *ShaderConfig) Definitions() string {
	lines := make([]string,len(config.definitions))
	lineIndex := 0
	for name,value := range config.definitions {
		lines[lineIndex] =
			fmt.Sprintf("#define %s %s",name,value)
	}
	return strings.Join(lines,"\n")+"\n"
}

func (config *ShaderConfig) compileShader(
		source string, glShaderType uint32,
) Shader {
	glShader := gl.CreateShader(glShaderType)
	csources, free :=
		gl.Strs(header+config.Definitions()+source+"\x00")
	defer free()
	gl.ShaderSource(glShader, 1, csources, nil)
	gl.CompileShader(glShader)
	shader := Shader(glShader)
	shader.checkCompileErrors()
	return shader
}

func (config *ShaderConfig) Compile() Program {
	vertexSource,err := ioutil.ReadFile(config.vertexPath)
	if err!=nil { log.Fatal(err) }
	fragmentSource,err := ioutil.ReadFile(config.fragmentPath)
	if err!=nil { log.Fatal(err) }

	vertexShader :=
		config.compileShader(string(vertexSource), gl.VERTEX_SHADER)
	defer vertexShader.Delete()

	fragmentShader :=
		config.compileShader(string(fragmentSource), gl.FRAGMENT_SHADER)
	defer fragmentShader.Delete()

	program := NewProgram()
	program.Attach(vertexShader)
	program.Attach(fragmentShader)
	program.Link()
	program.checkLinkErrors()
	return program
}

func NewProgram() Program {
	return Program(gl.CreateProgram())
}

func (p Program) Attach(s Shader) {
	gl.AttachShader(p.gl(), s.gl())
}

func (p Program) Link() {
	gl.LinkProgram(p.gl())
}

func (p Program) Use() {
	gl.UseProgram(uint32(p))
}

func (p Program) StopUsing() {
	gl.UseProgram(0)
}

func (p Program) SetBool(name string, value bool) {
	intval := 0
	if value {
		intval = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), int32(intval))
}
func (p Program) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), value)
}
func (p Program) SetUint(name string, value uint32) {
	gl.Uniform1ui(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), value)
}
func (p Program) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), value)
}
func (p Program) SetVec2(name string, value *mgl32.Vec2) {
	gl.Uniform2fv(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), 1, &value[0])
}
func (p Program) SetVec2F(name string, x, y float32) {
	gl.Uniform2f(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), x, y)
}

func (p Program) SetVec3(name string, value *mgl32.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), 1, &value[0])
}
func (p Program) SetVec3F(name string, x, y, z float32) {
	gl.Uniform3f(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), x, y, z)
}

func (p Program) SetVec4(name string, value *mgl32.Vec4) {
	gl.Uniform4fv(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), 1, &value[0])
}

func (p Program) SetVec4F(name string, x, y, z, w float32) {
	gl.Uniform4f(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), x, y, z, w)
}

func (p Program) SetMat2(name string, value *mgl32.Mat2) {
	gl.UniformMatrix2fv(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), 1, false, &value[0])
}

func (p Program) SetMat3(name string, value *mgl32.Mat3) {
	gl.UniformMatrix3fv(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), 1, false, &value[0])
}
func (p Program) SetMat4(name string, value *mgl32.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00")), 1, false, &value[0])
}

func (p Program) Locate(name string) int32 {
	return gl.GetUniformLocation(uint32(p), gl.Str(name+"\x00"))
}

func (p Program) SetMat3s(name string, values []mgl32.Mat3) {
	location := p.Locate(name)
	count := int32(len(values))
	gl.UniformMatrix3fv(location, count, false, &values[0][0])
}

func (p Program) SetMat4s(name string, values []mgl32.Mat4) {
	location := p.Locate(name)
	count := int32(len(values))
	gl.UniformMatrix4fv(location, count, false, &values[0][0])
}

func (p Program) SetVec2s(name string, values []mgl32.Vec2) {
	location := p.Locate(name)
	count := int32(len(values))
	gl.Uniform2fv(location, count, &values[0][0])
}

func (s Shader) GetInt(name uint32) int32 {
	var x int32
	gl.GetShaderiv(s.gl(), name, &x)
	return x
}

func (p Program) GetInt(name uint32) int32 {
	var x int32
	gl.GetProgramiv(p.gl(), name, &x)
	return x
}


func (s Shader) CompiledSuccessfully() bool {
	return s.GetInt(gl.COMPILE_STATUS) == gl.TRUE
}

func getShaderFlag(shader uint32, flagname uint32) bool {
	var flag int32
	gl.GetShaderiv(shader, flagname, &flag)
	return flag == gl.TRUE
}

func (s Shader) GetInfoLog() string {
	length := s.GetInt(gl.INFO_LOG_LENGTH)
	// add 1 to length for null terminator
	infolog := strings.Repeat("\x00", int(length)+1)
	gl.GetShaderInfoLog(s.gl(), length, nil, gl.Str(infolog))
	return infolog
}

func (p Program) GetInfoLog() string {
	length := p.GetInt(gl.INFO_LOG_LENGTH)
	// add 1 to length for null terminator
	infolog := strings.Repeat("\x00", int(length)+1)
	gl.GetProgramInfoLog(p.gl(), length, nil, gl.Str(infolog))
	return infolog
}

func (s Shader) checkCompileErrors() {
	if !s.CompiledSuccessfully() {
		log.Print("compile error")
		log.Fatal(s.GetInfoLog())
	}
}

func (p Program) LinkedSuccessfully() bool {
	return p.GetInt(gl.LINK_STATUS) == gl.TRUE
}

func (p Program) checkLinkErrors() {
	if !p.LinkedSuccessfully() {
		log.Print("link error")
		log.Fatal(p.GetInfoLog())
	}
}
