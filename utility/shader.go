package utility
// Helper shader program class from learnopengl.com

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
	ID uint32
}

func NewShader(vertexPath, fragmentPath string) *Shader {
	shader := &Shader{}

	vertexSource, err := ioutil.ReadFile(vertexPath)
	if err != nil {
		log.Fatal(err)
	}
	fragmentSource, err := ioutil.ReadFile(fragmentPath)
	if err != nil {
		log.Fatal(err)
	}

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(string(vertexSource) + "\x00")
	gl.ShaderSource(vertexShader, 1, csources, nil)
	free()
	gl.CompileShader(vertexShader)
	checkCompileErrors(vertexShader, "VERTEX")

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free = gl.Strs(string(fragmentSource) + "\x00")
	gl.ShaderSource(fragmentShader, 1, csources, nil)
	free()
	gl.CompileShader(fragmentShader)
	checkCompileErrors(fragmentShader, "FRAGMENT")

	shader.ID = gl.CreateProgram()
	gl.AttachShader(shader.ID, vertexShader)
	gl.AttachShader(shader.ID, fragmentShader)
	gl.LinkProgram(shader.ID)
	checkCompileErrors(shader.ID, "PROGRAM")

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shader
}

func (s *Shader) Use() {
	gl.UseProgram(s.ID)
}

func (s *Shader) SetBool(name string, value bool) {
	intval := 0
	if value {
		intval = 1
	}
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), int32(intval))
}
func (s *Shader) SetInt(name string, value int32) {
	gl.Uniform1i(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}
func (s *Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), value)
}
func (s *Shader) SetVec2(name string, value *mgl32.Vec2) {
	gl.Uniform2fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}
func (s *Shader) SetVec2F(name string, x, y float32) {
	gl.Uniform2f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), x, y)
}

func (s *Shader) SetVec3(name string, value *mgl32.Vec3) {
	gl.Uniform3fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}
func (s *Shader) SetVec3F(name string, x, y, z float32) {
	gl.Uniform3f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), x, y, z)
}

func (s *Shader) SetVec4(name string, value *mgl32.Vec4) {
	gl.Uniform4fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, &value[0])
}

func (s *Shader) SetVec4F(name string, x, y, z, w float32) {
	gl.Uniform4f(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), x, y, z, w)
}

func (s *Shader) SetMat2(name string, value *mgl32.Mat2) {
	gl.UniformMatrix2fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, false, &value[0])
}

func (s *Shader) SetMat3(name string, value *mgl32.Mat3) {
	gl.UniformMatrix3fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, false, &value[0])
}
func (s *Shader) SetMat4(name string, value *mgl32.Mat4) {
	gl.UniformMatrix4fv(gl.GetUniformLocation(s.ID, gl.Str(name+"\x00")), 1, false, &value[0])
}

func checkCompileErrors(shader uint32, tType string) {
	var success int32
	var infolog string
	if tType != "PROGRAM" {
		gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
		if success != 1 {
			var logInfoLength int32
			gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logInfoLength)

			infolog = strings.Repeat("\x00", int(logInfoLength)+1)
			gl.GetShaderInfoLog(shader, logInfoLength, nil, gl.Str(infolog))
			log.Fatal(infolog)
		}
	} else {
		gl.GetProgramiv(shader, gl.LINK_STATUS, &success)
		if success != 1 {
			infolog = strings.Repeat(" ", 1024)
			gl.GetProgramInfoLog(shader, 1024, nil, gl.Str(infolog))
			infolog = strings.Trim(infolog, " ")
			log.Fatal(infolog)
		}
	}
}
