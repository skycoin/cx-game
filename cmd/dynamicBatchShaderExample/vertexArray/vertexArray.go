package vertexArray

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBuffer"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexbufferLayout"
)

type vertexArray struct {
}

func (va *vertexArray) AddBuffer(vb *vertexbuffer.VertexBuffer, vbl *vertexbufferLayout.VertexbufferLayout) {

}

func RunVertxArray(indices []int, count int) *IndexBuffer {
	IB := &IndexBuffer{}
	IB.M_count = count
	gl.GenBuffers(1, &IB.M_renderID)
	//	fmt.Println("this is the ID Buffer: ", IB.M_renderID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, IB.M_renderID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, count*4, gl.Ptr(indices), gl.STATIC_DRAW)

	return IB
}

func AddBuffer(vb *vertexbuffer.VertexBuffer)

func (IB *IndexBuffer) DeleteBuffer() {
	gl.DeleteBuffers(1, &IB.M_renderID)
}

func (IB *IndexBuffer) Bind() {
	//	fmt.Println("this is the ID Buffer: ", IB.M_renderID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, IB.M_renderID)
}

func Unbind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}
