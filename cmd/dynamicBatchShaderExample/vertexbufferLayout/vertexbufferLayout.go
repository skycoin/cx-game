package vertexbufferLayout

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type VertexbufferElements struct {
	Type       uint32
	Count      uint32
	Normalized bool
}
type vertexbufferLayout struct {
	m_Elements []VertexbufferElements
	m_Stride   uint32
}

func (vbl *vertexbufferLayout) Push(receiveType interface{}, count uint32) {
	data := VertexbufferElements{}
	switch receiveType {
	case gl.INT:
		data = VertexbufferElements{
			gl.INT,
			count,
			false,
		}

	case gl.FLOAT:
		data = VertexbufferElements{
			gl.FLOAT,
			count,
			false,
		}

	case gl.UNSIGNED_INT:
		data = VertexbufferElements{
			gl.UNSIGNED_INT,
			count,
			false,
		}
	case gl.UNSIGNED_BYTE:
		data = VertexbufferElements{
			gl.UNSIGNED_INT,
			count,
			true,
		}
	default:
		// And here I'm feeling dumb. ;)
		fmt.Printf("type not supported yet.")
	}

	vbl.m_Elements = append(vbl.m_Elements, data)
	vbl.m_Stride += 4
}
