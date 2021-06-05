package particles

import (
	"github.com/skycoin/cx-game/render"
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

func InitParticles() {
	InitLasers()
}

func TickParticles(dt float32) {
	TickLasers(dt)
}

func DrawParticles(ctx render.Context) {
	DrawLasers(ctx)
}
