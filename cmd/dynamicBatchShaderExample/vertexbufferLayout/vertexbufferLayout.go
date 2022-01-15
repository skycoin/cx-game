package vertexbufferLayout

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexbufferElements struct {
	Type       uint32
	TypeSize   uint32
	Count      uint32
	Normalized bool
}
type VertexbufferLayout struct {
	m_Elements []VertexbufferElements
	m_Stride   uint32
}

func (vbl *VertexbufferLayout) Push(receiveType interface{}, count uint32) {

	data := VertexbufferElements{}
	var typeSize int32
	switch receiveType {
	case gl.INT:
		typeSize = 8
		data = VertexbufferElements{
			gl.INT,
			uint32(typeSize),
			count,
			false,
		}

	case gl.FLOAT:
		typeSize = 4
		data = VertexbufferElements{
			gl.FLOAT,
			uint32(typeSize),
			count,
			false,
		}

	case gl.UNSIGNED_INT:
		typeSize = 4
		data = VertexbufferElements{
			gl.UNSIGNED_INT,
			uint32(typeSize),
			count,
			false,
		}
	case gl.UNSIGNED_BYTE:
		typeSize = 1
		data = VertexbufferElements{
			gl.UNSIGNED_INT,
			uint32(typeSize),
			count,
			true,
		}
	default:
		// And here I'm feeling dumb. ;)
		fmt.Printf("type not supported yet.")
	}

	vbl.m_Elements = append(vbl.m_Elements, data)
	vbl.m_Stride += uint32(typeSize) * count
}

func (vbl *VertexbufferLayout) GetElements() []VertexbufferElements {
	return vbl.m_Elements
}

func (vbl *VertexbufferLayout) GetStride() uint32 {
	return vbl.m_Stride
}
