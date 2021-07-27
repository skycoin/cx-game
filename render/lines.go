package render

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var lineProgram Program

func (window *Window) DrawLines(
		lineArray []float32, color mgl32.Vec3, ctx Context,
) {
	// DEBUG: check if the array have the right amount of elements
	if len(lineArray) < 6 {
		log.Panicln("line array doesn't enough points to draw a line")
	} else if len(lineArray)%3 != 0 {
		log.Panicln("line array doesn't have the right amount of floats values to draw the lines")
	}

	lineProgram.Use()
	defer lineProgram.StopUsing()

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

	lineProgram.SetVec3("uColor",&color)

	mvp := ctx.MVP()

	lineProgram.SetMat4("uProjection",&mvp)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.LINES, 0, int32(len(lineArray)/2))
}
