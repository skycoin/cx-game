package ui

import (
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/cxmath/math32"
)

const arcTriangles = 50
var arcVAO uint32
var arcVBO uint32
var arcShader *utility.Shader

func createArcVertexAttributes() []float32 {
	radius := float32(0.5)
	attributes := make([]float32,arcTriangles*3*5)
	i := 0
	for tri := 0 ; tri < arcTriangles; tri++ {
		// for center point, set u=0.5 and v=0.5
		attributes[i+3] = 0.5
		attributes[i+4] = 0.5

		i += 5
		angle := 2 * math.Pi * float32(tri) / float32(arcTriangles)
		x := radius * -math32.Sin(angle)
		y := radius * math32.Cos(angle)
		z := float32(0)
		attributes[i] = x
		attributes[i+1] = y
		attributes[i+2] = z
		attributes[i+3] = 0.5+x
		attributes[i+4] = 0.5-y
		// arc is currently untextured so values of u and v are not important
		i += 5

		angle = 2 * math.Pi * float32(tri+1) / float32(arcTriangles)
		x = radius * -math32.Sin(angle)
		y = radius * math32.Cos(angle)
		z = float32(0)
		attributes[i] = x
		attributes[i+1] = y
		attributes[i+2] = z
		attributes[i+3] = 0.5+x
		attributes[i+4] = 0.5-y

		i += 5
	}
	return attributes
}

func initArcVAO() {
	var arcVBO uint32
	gl.GenBuffers(1,&arcVBO)
	gl.GenVertexArrays(1,&arcVAO)
	gl.BindVertexArray(arcVAO)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER,arcVBO)

	vertexAttributes := createArcVertexAttributes()

	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(vertexAttributes),
		gl.Ptr(vertexAttributes),
		gl.STATIC_DRAW,
	)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)
	//unbind
	gl.BindVertexArray(0)
}

var arcSprite spriteloader.SpriteID
func InitArc() {
	initArcVAO()
	arcShader = utility.NewShader(
		"./assets/shader/mvp.vert", "./assets/shader/arc.frag" )
	arcSprite = spriteloader.
		LoadSingleSprite("./assets/hud/hud_status_fill.png","hud_status_fill")
}

func DrawArc(mvp mgl32.Mat4, fullness float32) {
	gl.ActiveTexture(gl.TEXTURE0)
	arcShader.Use()
	defer arcShader.StopUsing()
	metadata := spriteloader.GetSpriteMetadata(arcSprite)
	arcShader.SetUint("tex", metadata.GpuTex)
	gl.BindTexture(gl.TEXTURE_2D, metadata.GpuTex)

	arcShader.SetMat4("mvp",&mvp)
	arcShader.SetVec4F("color",1,1,1,1)
	// FIXME note that this only works because sprite is in its own image.
	// Update this with a more robust strategy at some point.
	arcShader.SetVec2F("offset",0,0)
	arcShader.SetVec2F("scale",1,1)

	arcShader.SetFloat("value",2*math.Pi*fullness)

	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(spriteloader.QuadVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
