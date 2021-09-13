package particle_physics

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/cxmath"
)

var (
	knockbackVelocity = mgl32.Vec2{0,6}
	knockbackPower float32 = 0.05
)

func ApplyKnockback(target *physics.Body, bulletVelocity mgl32.Vec2) {
	extraVelocity := knockbackVelocity.Add(bulletVelocity.Mul(knockbackPower))

	target.Vel =
		target.Vel.Add(cxmath.Vec2{extraVelocity.X(),extraVelocity.Y()})
}
