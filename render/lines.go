package render

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var lineProgram uint32

func (window *Window) DrawLine(x0, y0, z0, x1, y1, z1 float32) {
	var lineArray = []float32{
		x0, y0, z0,
		x1, y1, z1,
	}
	color := []float32{1.0, 1.0, 1.0}

	window.DrawLines(lineArray, color)
}

func (window *Window) DrawLines(lineArray []float32, color []float32) {
	// DEBUG: check if the array have the right amount of elements
	if len(lineArray) < 6 {
		log.Panicln("line array doesn't enough points to draw a line")
	} else if len(lineArray)%3 != 0 {
		log.Panicln("line array doesn't have the right amount of floats values to draw the lines")
	}
	if len(color) > 4 || len(color) < 3 {
		log.Panicln("wrong amount of floats values for a color (need 3 or 4)")
	}

	gl.UseProgram(lineProgram)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	var vao uint32
	gl.GenVertexArrays(1, &vao)

	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(lineArray),
		gl.Ptr(lineArray),
		gl.STATIC_DRAW,
	)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.Uniform3fv(
		gl.GetUniformLocation(lineProgram, gl.Str("uColor\x00")),
		1,
		&color[0],
	)

	// unforms
	worldTransform := mgl32.Mat4.Mul4(
		mgl32.Translate3D(float32(0.0), float32(0.0), -10),
		mgl32.Scale3D(float32(1.0), float32(1.0), 1),
	)

	aspect := float32(window.Width) / float32(window.Height)
	projectTransform := mgl32.Perspective(
		mgl32.DegToRad(45), aspect, 0.1, 100.0,
	)

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(lineProgram, gl.Str("uWorld\x00")),
		1, false, &worldTransform[0],
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(lineProgram, gl.Str("uProjection\x00")),
		1, false, &projectTransform[0],
	)

	gl.BindVertexArray(vao)
	//gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.DrawArrays(gl.LINES, 0, int32(len(lineArray)/2))
}
