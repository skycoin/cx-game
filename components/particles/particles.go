package particles

import "github.com/skycoin/cx-game/cxmath"

// a particle that bounces but is solid

// a particle that does not bounce but is alpha blended and not solid

// a particle that floats and is alpha blended (drifts in air, no gravity)

// And two draw modes

// 1> solid particle

// 2> alpha blended particle (has transparency, glowing, etc)

// And three physics modes

// 1> bounces, gravity

// 2> Does not bounce, stops on impact, gravity

// 3> "drifts" at fixed velocity, no gravity

// when you fire gun do "shell casing" particle coming out of player
// when you hit block, do "debris" particle, etc
// test a spark or type of glowing/alpha blended particle
// create a struct called "Emitter" which is just something that creates particles, we can throw in game

type Particle struct {
	Position    cxmath.Vec2
	Velocity    cxmath.Vec2
	TimeToLive  float32
	DrawMode    DrawMode
	PhysicsMode PhysicsMode
}

type DrawMode uint32

const (
	SolidDrawMode DrawMode = iota
	AlphaBlendedDrawMode
)

type PhysicsMode uint32

const (
	//bounces, gravity
	BouncePhysicsMode PhysicsMode = iota
	//does not bounce, gravity
	NoBouncePhysicsMode
	//
	DriftPhysicsMode
)

type ParticleID uint32

type Emitter struct {
	ParticleID ParticleID
}

func (emitter *Emitter) Emit() {}
