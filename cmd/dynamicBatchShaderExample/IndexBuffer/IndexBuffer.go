package indexbuffer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type IndexBuffer struct {
	M_renderID uint32
	M_count    int
}

func RunIndexBuffer(indices []uint32, count int) *IndexBuffer {
	IB := &IndexBuffer{}
	IB.M_count = count
	gl.GenBuffers(1, &IB.M_renderID)
	fmt.Println("this is the ID Buffer: ", IB.M_count)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, IB.M_renderID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, count*4, gl.Ptr(indices), gl.STATIC_DRAW)

	return IB
}

func (IB *IndexBuffer) DeleteBuffer() {
	gl.DeleteBuffers(1, &IB.M_renderID)
}

func (IB *IndexBuffer) Bind() {
	//	fmt.Println("this is the ID Buffer: ", IB.M_renderID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, IB.M_renderID)
}

func (IB *IndexBuffer) Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}
