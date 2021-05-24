package render

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var lineProgram uint32

func (window *Window) DrawLines(
		lineArray []float32, color []float32, ctx Context,
) {
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

	mvp := ctx.MVP()

	gl.UniformMatrix4fv(
		gl.GetUniformLocation(lineProgram, gl.Str("uProjection\x00")),
		1, false, &mvp[0],
	)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.LINES, 0, int32(len(lineArray)/2))
}
