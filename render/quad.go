package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

var QuadVao uint32

func InitQuadVao() {
	QuadVao = MakeQuadVao()
}

// x,y,z,u,v
var quadVertexAttributes = []float32{
	0.5, 0.5, 0, 1, 0,
	0.5, -0.5, 0, 1, 1,
	-0.5, 0.5, 0, 0, 0,

	0.5, -0.5, 0, 1, 1,
	-0.5, -0.5, 0, 0, 1,
	-0.5, 0.5, 0, 0, 0,
}


func MakeQuadVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(quadVertexAttributes),
		gl.Ptr(quadVertexAttributes),
		gl.STATIC_DRAW,
	)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)
	//unbind
	gl.BindVertexArray(0)

	return vao
}
