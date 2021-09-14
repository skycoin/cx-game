package render

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	lineProgram Program
	lines_vao   uint32
	lines_vbo   uint32
)

func InitDrawLines() {
	gl.GenVertexArrays(1, &lines_vao)
	gl.BindVertexArray(lines_vao)

	gl.GenBuffers(1, &lines_vbo)

	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, lines_vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*4,
		nil,
		gl.STATIC_DRAW,
	)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
}

func DrawLines(
	lineArray []float32, color mgl32.Vec3, mvp mgl32.Mat4,
) {
	// DEBUG: check if the array have the right amount of elements
	if len(lineArray) < 4 {
		return
		// log.Panicln("line array doesn't enough points to draw a line")
	} else if len(lineArray)%2 != 0 {
		log.Panicln("line array doesn't have the right amount of floats values to draw the lines")
	}

	lineProgram.Use()
	gl.BindBuffer(gl.ARRAY_BUFFER, lines_vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(lineArray)*4,
		gl.Ptr(lineArray),
		gl.STATIC_DRAW,
	)

	defer lineProgram.StopUsing()
	lineProgram.SetVec3("uColor", &color)
	lineProgram.SetMat4("uProjection", &mvp)

	gl.BindVertexArray(lines_vao)
	gl.DrawArrays(gl.LINES, 0, int32(len(lineArray)/2))
}
