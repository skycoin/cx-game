package spineloader

import (
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/renderer"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/shader"

	// "github.com/skycoin/cx-game/render"
	"github.com/go-gl/gl/v4.1-core/gl"
	indexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBufferDY"
	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBufferDY"
	vertexArray "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexArrayDY"
	vertexbufferLayout "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexbufferLayoutDY"
)

var Shader *shader.Shader
var Render *renderer.Render

var Ib *indexbuffer.IndexBuffer
var Vb *vertexbuffer.VertexBuffer
var Va *vertexArray.VertexArray
var Vbl *vertexbufferLayout.VertexbufferLayout

var TriIb *indexbuffer.IndexBuffer
var TriVb *vertexbuffer.VertexBuffer
var TriVa *vertexArray.VertexArray
var TriVbl *vertexbufferLayout.VertexbufferLayout

var LineIb *indexbuffer.IndexBuffer
var LineVb *vertexbuffer.VertexBuffer
var LineVa *vertexArray.VertexArray
var LineVbl *vertexbufferLayout.VertexbufferLayout

const MaxQuadCount = 1000
const MaxVertexCount = MaxQuadCount * 4
const MaxIndexCount = MaxQuadCount * 6

const MaxTriangesCount = 1000
const MaxTriangleVertexCount = MaxTriangesCount * 3
const MaxTriangleIndexCount = MaxTriangesCount * 3

const MaxLineCount = 1000
const MaxLineVertexCount = MaxLineCount * 2
const MaxLineIndexCount = MaxLineCount * 2

var Indices = make([]uint32, MaxIndexCount)

var TriangleIndices = make([]uint32, MaxTriangleIndexCount)
var LineIndices = make([]uint32, MaxLineIndexCount)

func InitSpineSpriteLoader() {
	// Program = render.CompileProgram(
	// 	"./assets/shader/spine/spine.vert", "./assets/shader/spine/spine.frag",
	// )
	Shader = shader.SetupShader("./assets/shader/spine/cxSpine.shader")
	Render = renderer.SetupRender()
	SetupInices()
	SetupBuffers()

}
func SetupBuffers() {
	//setup vertex array
	Va = vertexArray.SetUpVertxArray()
	// setup and run vertex buffer
	Vb = vertexbuffer.RunVertexBuffer(5 * 4 * MaxVertexCount)
	//setup vertex layout
	Vbl = &vertexbufferLayout.VertexbufferLayout{}
	Vbl.Push(gl.FLOAT, 2)
	Vbl.Push(gl.FLOAT, 2)
	Vbl.Push(gl.FLOAT, 1)
	Va.AddBuffer(Vb, Vbl)

	// setup and run index buffer
	Ib = indexbuffer.RunIndexBuffer(Indices, len(Indices))

	//setup vertex array
	TriVa = vertexArray.SetUpVertxArray()
	// setup and run vertex buffer
	TriVb = vertexbuffer.RunVertexBuffer(5 * 4 * MaxVertexCount)
	//setup vertex layout
	TriVbl = &vertexbufferLayout.VertexbufferLayout{}
	TriVbl.Push(gl.FLOAT, 2)
	TriVbl.Push(gl.FLOAT, 2)
	TriVbl.Push(gl.FLOAT, 1)
	TriVa.AddBuffer(TriVb, TriVbl)

	// setup and run index buffer
	TriIb = indexbuffer.RunIndexBuffer(LineIndices, len(LineIndices))

	//setup vertex array
	LineVa = vertexArray.SetUpVertxArray()
	// setup and run vertex buffer
	LineVb = vertexbuffer.RunVertexBuffer(5 * 4 * MaxVertexCount)
	//setup vertex layout
	LineVbl = &vertexbufferLayout.VertexbufferLayout{}

	//add vertex buffer to vertex bufferlayout

	LineVbl.Push(gl.FLOAT, 2)
	LineVbl.Push(gl.FLOAT, 2)
	LineVbl.Push(gl.FLOAT, 1)
	LineVa.AddBuffer(LineVb, LineVbl)

	// setup and run index buffer
	LineIb = indexbuffer.RunIndexBuffer(TriangleIndices, len(TriangleIndices))
}

func SetupInices() {
	var offset uint32 = 0
	for i := 0; i < MaxIndexCount; i += 6 {
		Indices[i+0] = 0 + offset
		Indices[i+1] = 1 + offset
		Indices[i+2] = 2 + offset

		Indices[i+3] = 2 + offset
		Indices[i+4] = 3 + offset
		Indices[i+5] = 0 + offset

		offset += 4
	}

	offset = 0
	for i := 0; i < MaxTriangleIndexCount; i += 3 {
		TriangleIndices[i+0] = 0 + offset
		TriangleIndices[i+1] = 1 + offset
		TriangleIndices[i+2] = 2 + offset

		offset += 3
	}

	offset = 0
	for i := 0; i < MaxLineIndexCount; i += 2 {
		LineIndices[i+0] = 0 + offset
		LineIndices[i+1] = 1 + offset

		offset += 2
	}
}
