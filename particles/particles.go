package particles

import (
	"log"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/spriteloader"
)

type Particle struct {
	ID               int32
	ParticleMetaType int32
	Type             int32
	Sprite           uint32
	Size             int32
	PosX             float32
	PosY             float32
	VelocityX        float32
	VelocityY        float32
	TimeToLive       float32
}
var particles = []Particle{}
var particleShader *utility.Shader
const initialVelocityScale = 3
const tileChunkLifetime = 1

func InitParticles() {
	particleShader = utility.NewShader(
		"./assets/shader/simple.vert", "./assets/shader/particle.frag" )
	InitLasers()
}

func TickParticles(dt float32) {
	TickLasers(dt)
	// age and kill particles
	newParticles := []Particle{}
	for _,laser := range particles {
		laser.TimeToLive -= dt
		if laser.TimeToLive > 0 {
			newParticles = append(newParticles,laser)
		}
	}
	particles = newParticles
}

func configureGlForParticles() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_2D)
	// bits of tile are still pixel art; use nearest interpolation
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.BindVertexArray(spriteloader.QuadVao);
}

func DrawParticles(ctx render.Context) {
	DrawLasers(ctx)
	configureGlForParticles()
	for _,particle := range particles {
		DrawParticle(particle,ctx)
	}
}

func DrawParticle(particle Particle, ctx render.Context) {
	log.Print("drawing particle")
}

func CreateTileChunk(x,y float32, TileSpriteID uint32) {
	particle := Particle {
		ID: rand.Int31(),
		PosX: x, PosY: y,
		VelocityX: (rand.Float32()-0.5)*initialVelocityScale,
		VelocityY: (rand.Float32()-0.5)*initialVelocityScale,
		Sprite: TileSpriteID,
		TimeToLive: tileChunkLifetime,
	}
	particles = append(particles, particle)
}
