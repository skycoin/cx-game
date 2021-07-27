package particles

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

type Laser struct {
	ttl       float32    // time to live
	length    float32    // sprite segments required to draw the laser
	transform mgl32.Mat4 // transform relative to world
}

type texture struct {
	gpuTex        uint32 // identifier for OpenGL
	width, height int
	vao           uint32
}

var laserTex texture

const laserDuration = 1.0 // time in seconds that a laser lasts
var lasers []Laser = []Laser{}
var laserProgram render.Program

const segmentLength = 0.8

// number of times the laser texture animates throughout its life
const laserAnimSpeed = 6

func InitLasers() {
	laserProgram = render.CompileProgram(
		"./assets/shader/mvp.vert", "./assets/shader/laser.frag")

	_, img, _ :=
		spriteloader.LoadPng("./assets/projectile/laser_beam_core_00.png")

	laserTex.width = img.Bounds().Dx()
	laserTex.height = img.Bounds().Dy()
	laserTex.gpuTex = spriteloader.MakeTexture(img)
	laserTex.vao = makeLaserVao()
}

// quad where (0,0) is vertically centered, but horizontally aligned to left
func makeLaserVao() uint32 {
	// x,y,z,u,v
	var quadVertexAttributes = []float32{
		1, 0.5, 0, 1, 0,
		1, -0.5, 0, 1, 1,
		0, 0.5, 0, 0, 0,

		1, -0.5, 0, 1, 1,
		0, -0.5, 0, 0, 1,
		0, 0.5, 0, 0, 0,
	}
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(quadVertexAttributes),
		gl.Ptr(quadVertexAttributes),
		gl.STATIC_DRAW,
	)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)

	return vao
}
func CreateLaser(from, to mgl32.Vec2) {
	disp := to.Sub(from)
	right := mgl32.Vec2{1, 0}
	angle := -cxmath.AngleTo(right, disp)
	length := disp.Len()
	transform := mgl32.Ident4().
		Mul4(mgl32.Translate3D(from.X(), from.Y(), 0)).
		Mul4(mgl32.HomogRotate3DZ(angle)).
		Mul4(mgl32.Scale3D(length, 1, 1))

	laser := Laser{
		transform: transform,
		ttl:       laserDuration,
		length:    length,
	}
	lasers = append(lasers, laser)
}

func TickLasers(dt float32) {
	newLasers := []Laser{}
	for _, laser := range lasers {
		laser.ttl -= dt
		if laser.ttl > 0 {
			newLasers = append(newLasers, laser)
		}
	}
	lasers = newLasers
}

func DrawLaser(laser Laser, ctx render.WorldContext) {
	alpha := laser.ttl / laserDuration
	laserProgram.SetVec4F("color", 1, 1, 1, alpha)

	mvp := ctx.ModelToModelViewProjection(laser.transform)
	laserProgram.SetMat4("mvp", &mvp)

	laserProgram.SetVec2F("offset", alpha*laserAnimSpeed, 0)
	laserProgram.SetVec2F("scale", laser.length, 1)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func configureGlForLaser() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_2D)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, laserTex.gpuTex)
	// blurry is better than jagged for a laser
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.BindVertexArray(laserTex.vao)
}

func DrawLasers(ctx render.WorldContext) {
	// only need to set up shader once for all lasers
	laserProgram.Use()
	projection := ctx.Projection()
	laserProgram.SetMat4("projection", &projection)
	laserProgram.SetUint("tex", laserTex.gpuTex)
	configureGlForLaser()

	for _, laser := range lasers {
		DrawLaser(laser, ctx)
	}
}
