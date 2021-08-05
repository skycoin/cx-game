package particle_emitter

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/cxmath"
)

func CreateBullet(position, velocity mgl32.Vec2) {
	bulletEmitter.Emit(
		cxmath.Vec2{position.X(), position.Y()},
		cxmath.Vec2{velocity.X(), velocity.Y()},
	)
}
