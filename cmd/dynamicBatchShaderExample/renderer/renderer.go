package renderer

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	indexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBuffer"
	indexbufferDY "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBufferDY"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/shader"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexArray"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexArrayDY"
)

type GL_ERROR_CODE int32

const (
	NONE GL_ERROR_CODE = 0
)

type Render struct {
}

func SetupRender() *Render {
	render := &Render{}
	return render
}

func GLClearError() {
	for gl.GetError() != gl.NO_ERROR {
		fmt.Println("Cleared an Error")
	}
}

func GLCheckError() {
	for err := gl.GetError(); err != 0; {
		fmt.Println("got a An Error: ", err)
	}
}

func (R *Render) Draw(va *vertexArray.VertexArray, ib *indexbuffer.IndexBuffer, shader *shader.Shader) {
	shader.Bind()
	va.Bind()
	ib.Bind()

	gl.DrawElements(gl.TRIANGLES, int32(ib.GetCount()), gl.UNSIGNED_INT, nil)
}
func (R *Render) DrawDY(va *vertexArrayDY.VertexArray, ib *indexbufferDY.IndexBuffer, shader *shader.Shader) {
	shader.Bind()
	va.Bind()
	ib.Bind()

	gl.DrawElements(gl.TRIANGLES, int32(ib.GetCount()), gl.UNSIGNED_INT, nil)
}

func (R *Render) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
