package particles

import (
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
)

type Particle struct {
	ID               int32
	ParticleMetaType int32
	Type             int32
	Sprite           int32
	Size             int32
	PosX             int32
	PosY             int32
	VelocityX        float32
	VelocityY        float32
	TimeToLive       int32
}
var particles = []Particle{}
var particleShader *utility.Shader

func InitParticles() {
	particleShader = utility.NewShader(
		"./assets/shader/simple.vert", "./assets/shader/particle.frag" )
	InitLasers()
}

func TickParticles(dt float32) {
	TickLasers(dt)
}

func DrawParticles(ctx render.Context) {
	DrawLasers(ctx)
}
