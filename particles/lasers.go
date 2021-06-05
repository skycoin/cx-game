package particles

// NOTE laser rendering can probably be optimized by
// creating a large VAO that renders 20-50 segments
// and then using a single draw call per laser

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/spriteloader"
)

type Laser struct {
	ttl float32 // time to live
	length float32 // sprite segments required to draw the laser
	transform mgl32.Mat4 // transform relative to world
}

type texture struct {
	gpuTex uint32 // identifier for OpenGL
	width, height int
	vao uint32
}
var laserTex texture

const laserDuration = 1.0 // time in seconds that a laser lasts
var lasers []Laser = []Laser{}
var laserShader *utility.Shader
const segmentLength = 0.8

func InitLasers() {
	laserShader = utility.NewShader(
		"./assets/shader/simple.vert", "./assets/shader/laser.frag" )

	_,img,_ :=
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
	right := mgl32.Vec2 { 1, 0 }
	angle := - cxmath.AngleTo(right,disp)
	transform := mgl32.Ident4().
		Mul4(mgl32.Translate3D(from.X(),from.Y(),0)).
		Mul4(mgl32.HomogRotate3DZ(angle))

	laser := Laser {
		transform: transform,
		ttl: laserDuration,
		length: disp.Len(),
	}
	lasers = append(lasers,laser)
}

func TickLasers(dt float32) {
	newLasers := []Laser{}
	for _,laser := range lasers {
		laser.ttl -= dt
		if laser.ttl > 0 {
			newLasers = append(newLasers,laser)
		}
	}
	lasers = newLasers
}

func DrawLaser(laser Laser, ctx render.Context) {
	laserShader.SetFloat("start",0)
	laserShader.SetFloat("stop",1)
	alpha := laser.ttl / laserDuration
	laserShader.SetVec4F("color", 1,1,1, alpha)

	lastX := float32(0)
	for x:=float32(0); x<laser.length-1; x+=segmentLength {
		world := ctx.World.
			Mul4(laser.transform).
			Mul4(mgl32.Translate3D(x,0,0))
		laserShader.SetMat4("world", &world)
		// 1 quad = 2 (3d) verts = 6 floats
		// texture has already been bound by configureGlForLaser()
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
		lastX = x
	}
	// last segment is special case; we want to draw EXACTLY up to the length
	x := lastX + segmentLength
	world := ctx.World.
		Mul4(laser.transform).
		Mul4(mgl32.Translate3D(x,0,0))
	laserShader.SetMat4("world", &world)

	laserShader.SetFloat("stop", laser.length - x)
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

	gl.BindVertexArray(laserTex.vao);
}

func DrawLasers(ctx render.Context) {
	// only need to set up shader once for all lasers
	laserShader.Use()
	laserShader.SetMat4("projection", &ctx.Projection)
	laserShader.SetUint("tex", laserTex.gpuTex)
	configureGlForLaser()

	for _,laser := range lasers {
		DrawLaser(laser,ctx)
	}
}
