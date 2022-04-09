package agent_draw

import (
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	indexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBufferDY"
	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBufferDY"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/character"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/renderer"
	vertexArray "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexArrayDY"
	vertexbufferLayout "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexbufferLayoutDY"
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/test/spine-animation/animation"
)

var ib *indexbuffer.IndexBuffer
var vb *vertexbuffer.VertexBuffer
var va *vertexArray.VertexArray
var vbl *vertexbufferLayout.VertexbufferLayout

var (
	characters     []*character.Character
	character_l    *character.Character
	characterIndex int

	characterRenderer *renderer.Render
)

var (
	//positions = []float32{}

	samplers = [31]int32{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
	}
)

const MaxQuadCount = 1000
const MaxVertexCount = MaxQuadCount * 4
const MaxIndexCount = MaxQuadCount * 6

var indices = make([]uint32, MaxIndexCount)
var offset uint32 = 0

var isSetUp = false

func startUp() {

	for _, loc := range animation.LoadList("./test/spine-animation/animation") {
		character, err := character.LoadCharacter(loc)
		if err != nil {
			log.Println(loc.Name, err)
			break
		}

		for _, skin := range character.Skeleton.Data.Skins {
			for _, att := range skin.Attachments {
				if _, ismesh := att.(*spine.MeshAttachment); ismesh {
					log.Println(loc.Name, "Unsupported")
					break
				}
			}
		}

		characters = append(characters, character)
	}

	for i := 0; i < MaxIndexCount; i += 6 {
		indices[i+0] = 0 + offset
		indices[i+1] = 1 + offset
		indices[i+2] = 2 + offset

		indices[i+3] = 2 + offset
		indices[i+4] = 3 + offset
		indices[i+5] = 0 + offset

		offset += 4
	}
	//setup vertex array
	va = vertexArray.SetUpVertxArray()
	// setup and run vertex buffer
	vb = vertexbuffer.RunVertexBuffer(5 * 4 * MaxVertexCount)
	//setup vertex layout
	vbl = &vertexbufferLayout.VertexbufferLayout{}
	//add vertex buffer to vertex bufferlayout

	vbl.Push(gl.FLOAT, 2)
	vbl.Push(gl.FLOAT, 2)
	vbl.Push(gl.FLOAT, 1)
	va.AddBuffer(vb, vbl)

	// setup and run index buffer
	ib = indexbuffer.RunIndexBuffer(indices, len(indices))
	character_l = characters[0]

	characterRenderer = renderer.SetupRender()
	isSetUp = true
	character.Shader.SetUniform1iv("u_Texture", int32(len(samplers)), &samplers[0])
}

func SpineAnimatedDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	character.Shader.Bind()
	defer character.Shader.UnBind()

	if isSetUp == false {
		startUp()
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// character.Shader.SetUniForm4f("u_Color", 0.8, 0.3, 0.8, 1.0)

	var positions []float32
	positions = append(positions, 0, 0, 0, 0, 0)
	// positions = character_l.Draw()
	vb.Bind()
	vb.BufferSubData(positions)

	for _, agent := range agents {
		// tex := agent.AnimationPlayback.Animation.Texture
		// gl.ActiveTexture(gl.TEXTURE0)
		// gl.BindTexture(gl.TEXTURE_2D, tex)

		alpha := timer.GetTimeBetweenTicks() / constants.PHYSICS_TICK

		var interpolatedPos cxmath.Vec2
		if !agent.Transform.PrevPos.Equal(agent.Transform.Pos) {
			interpolatedPos = agent.Transform.PrevPos.Mult(1 - alpha).Add(agent.Transform.Pos.Mult(alpha))

		} else {
			interpolatedPos = agent.Transform.Pos
		}
		fmt.Println(agent)
		character_l.Update(1/float64(alpha), float64(interpolatedPos.X), float64(interpolatedPos.Y))

		translate := mgl32.Translate3D(
			interpolatedPos.X,
			interpolatedPos.Y,
			constants.AGENT_Z,
		)
		scale := mgl32.Scale3D(
			agent.Transform.Size.X*agent.Transform.Direction,
			agent.Transform.Size.Y,
			1,
		)
		transform := translate.Mul4(scale)
		wrappedTransform := wrapSpineTransform(
			transform,
			ctx.Camera.PlanetWidth,
			ctx.Camera.GetTransform(),
		)
		projection := spriteloader.Window.GetProjectionMatrix()
		mvp := projection.Mul4(ctx.Camera.GetViewMatrix()).Mul4(wrappedTransform)

		character.Shader.Bind()
		character.Shader.SetUniFormMat4f("u_MVP", mvp)
		characterRenderer.DrawDY(va, ib, character.Shader)

		// anim.Program.SetMat4("mvp", &mvp)
		// texTransform := agent.AnimationPlayback.Frame().Transform
		// anim.Program.SetMat3("texTransform", &texTransform)
		// gl.DrawArrays(gl.TRIANGLES, 0, 6)
	}
}

func wrapSpineTransform(raw mgl32.Mat4, worldWidth float32, cameraTransform mgl32.Mat4) mgl32.Mat4 {
	rawX := raw.At(0, 3)
	x := math32.PositiveModulo(rawX, worldWidth)
	camX := cameraTransform.At(0, 3)
	if x-camX > worldWidth/2 {
		x -= worldWidth
	}
	if x-camX < -worldWidth/2 {
		x += worldWidth
	}

	translate := mgl32.Translate3D(x-rawX, 0, 0)
	return translate.Mul4(raw)
}
