package vertexArrayDY

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBufferDY"
	vertexbufferLayout "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexbufferLayoutDY"
)

type VertexArray struct {
	M_renderID uint32
}

func SetUpVertxArray() *VertexArray {
	VA := &VertexArray{}
	gl.GenVertexArrays(1, &VA.M_renderID)

	return VA
}

func (va *VertexArray) AddBuffer(vb *vertexbuffer.VertexBuffer, vbl *vertexbufferLayout.VertexbufferLayout) {
	va.Bind()
	vb.Bind()
	var index uint32
	var offset uint32 = 0
	elements := vbl.GetElements()

	for index = 0; index < uint32(len(elements)); index++ {
		element := elements[index]
		gl.EnableVertexAttribArray(index)
		gl.VertexAttribPointer(index, int32(element.Count), element.Type, element.Normalized, int32(vbl.GetStride()), gl.PtrOffset(int(offset)))
		offset += element.Count * element.TypeSize
	}

}

func (VA *VertexArray) DeleteBuffer() {
	gl.DeleteVertexArrays(1, &VA.M_renderID)
}

func (VA *VertexArray) Bind() {
	//	fmt.Println("this is the ID Buffer: ", IB.M_renderID)
	gl.BindVertexArray(VA.M_renderID)
}

func (VA *VertexArray) Unbind() {
	gl.BindVertexArray(0)
}
