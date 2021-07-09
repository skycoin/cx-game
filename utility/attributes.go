package utility

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Geometry struct {
	data []float32
}

func NewGeometry() Geometry {
	return Geometry { data: []float32{} }
}

type Vert struct {
	X,Y,Z,U,V float32
}

func (g *Geometry) AddVert(v Vert) {
	log.Printf("adding Vert: %+v",v)
	g.data = append(g.data, v.X,v.Y,v.Z,v.U,v.V)
}

func (g *Geometry) Verts() int32 {
	return int32(len(g.data)/5)
}

func (g *Geometry) AddTri(a,b,c Vert) {
	g.AddVert(a)
	g.AddVert(b)
	g.AddVert(c)
}

func (g *Geometry) AddQuad(a,b,c,d Vert) {
	g.AddTri(a,b,c)
	g.AddTri(c,d,a)
}

func (g *Geometry) AddQuadFromCorners(from,to Vert) {
	a := from
	b := Vert { X: to.X, U: to.U, Y: from.Y, V: from.V }
	c := to
	d := Vert { X: from.X, U: from.U, Y: to.Y, V: to.V }
	g.AddQuad(a,b,c,d)
}

// upload this geometry to the GPU
func (g *Geometry) Upload() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(g.data), // 32-bit floats are 4 bytes long
		gl.Ptr(g.data),
		gl.STATIC_DRAW,
	)
	// communicate the (x,y,z,u,v) data layout to the GPU
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)
	//unbind
	gl.BindVertexArray(0)

	return vao
}
