package render

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	lineShader *Shader
	vao        uint32
	vbo        uint32
)

func InitDrawLines() {
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)

	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*4,
		nil,
		gl.STATIC_DRAW,
	)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
}

func (window *Window) DrawLines(
	lineArray []float32, color []float32, ctx Context,
) {
	// DEBUG: check if the array have the right amount of elements
	if len(lineArray) < 4 {
		log.Panicln("line array doesn't enough points to draw a line")
	} else if len(lineArray)%2 != 0 {
		log.Panicln("line array doesn't have the right amount of floats values to draw the lines")
	}
	if len(color) > 4 || len(color) < 3 {
		log.Panicln("wrong amount of floats values for a color (need 3 or 4)")
	}

	lineShader.Use()

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(lineArray),
		gl.Ptr(lineArray),
		gl.STATIC_DRAW,
	)

	lineShader.SetVec3F("uColor", color[0], color[1], color[2])
	mvp := ctx.MVP()
	lineShader.SetMat4("uProjection", &mvp)

	gl.DrawArrays(gl.LINES, 0, int32(len(lineArray)/2))
}
