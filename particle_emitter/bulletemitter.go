package particle_emitter

import (
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
)

var bulletEmitter *BulletEmitter

type BulletEmitter struct {
	particleList *particles.ParticleList
}

func NewBulletEmitter(particlelist *particles.ParticleList) *BulletEmitter {
	bulletEmitter := BulletEmitter{
		particleList: particlelist,
	}

	return &bulletEmitter
}

func (emitter *BulletEmitter) Emit(position, velocity cxmath.Vec2) {
	emitter.particleList.AddParticle(
		position,
		//set velocity up
		velocity,
		0.3,
		0.5,
		0.1,
		spriteloader.GetSpriteIdByNameUint32("star"),
		3,
		constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED,
		constants.PARTICLE_PHYSICS_HANDLER_DISSAPPEAR_ON_HIT_CALLBACK,
		emitter.OnHitCallback(),
	)
}

func (emitter *BulletEmitter) GetVelocity() cxmath.Vec2 {
	// velocity := cxmath.Vec2{}
	// velocity.X = rand.Float32() - 0.5*2
	// velocity.Y = -emitter.Vel.Y * rand.Float32()

	// return velocity
	// inputVec := input.GetMouseWorldCoords().Mgl32()
	return cxmath.Vec2{}
	// direction := inputVec.Sub(emitter.Pos.Mgl32()).Normalize()

	// result := direction.Mul(20)
	// return cxmath.Vec2{
	// 	result.X(),
	// 	result.Y(),
	// }
}

// func (emitter *BulletEmitter) GetDirection() cxmath.Vec2 {
// 	newDirection := mgl32.Vec2{input.GetMouseWorldCoordsX(), input.GetMouseWorldCoordsY()}
// 	angle := math.Acos(float64(
// 		emitter.Pos.Mgl32().Dot(newDirection) / (emitter.Pos.Length() * newDirection.Len()),
// 	))

// 	rotatedX := float32(math.Cos(float64(angle)))*emitter.Vel.X -
// 		float32(math.Sin(float64(angle)))*emitter.Vel.Y
// 	rotatedY := float32(math.Sin(float64(angle)))*emitter.Vel.X +
// 		float32(math.Cos(float64(angle)))*emitter.Vel.Y

// 	// fmt.Println("ROTATED: ", rotatedX, "   ", rotatedY, " SCREENX: ", input.GetMouseWorldCoordsX(), "  ", input.GetMouseWorldCoordsY())

// 	return cxmath.Vec2{rotatedX, rotatedY}
// }

func (emitter *BulletEmitter) OnHitCallback() func(*particles.Particle) {
	return func(particle *particles.Particle) {
		sparkEmitter.Emit(particle)
	}
}
