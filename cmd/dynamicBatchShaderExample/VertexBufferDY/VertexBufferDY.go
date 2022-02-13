package vertexbufferDY

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexBuffer struct {
	M_renderID uint32
}

func RunVertexBuffer(size int) *VertexBuffer {
	vb := &VertexBuffer{}
	gl.GenBuffers(1, &vb.M_renderID)
	gl.BindBuffer(gl.ARRAY_BUFFER, vb.M_renderID)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(nil), gl.DYNAMIC_DRAW)
	// m_renderID = m_rID
	//vb.Bind()
	return vb
}

func (VB *VertexBuffer) DeleteBuffer() {
	gl.DeleteBuffers(1, &VB.M_renderID)
}

func (VB *VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, VB.M_renderID)
}

func (VB *VertexBuffer) Unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
