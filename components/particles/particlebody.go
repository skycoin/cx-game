package particles

import (
	"math"

	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/world/worldcollider"
)

type ParticleBody struct {
	//todo add previous position and velocity if needed
	Pos     cxmath.Vec2
	Vel     cxmath.Vec2
	PrevPos cxmath.Vec2
	PrevVel cxmath.Vec2
	Size    cxmath.Vec2

	Collisions physics.CollisionInfo

	Elasticity float32
	Friction   float32
}

type bodyBounds struct {
	left, right, top, bottom float32
}

func round32(x float32) float32 {
	return float32(math.Round(float64(x)))
}

func (pb ParticleBody) bounds(newpos cxmath.Vec2) bodyBounds {
	left := round32(newpos.X - pb.Size.X/2)
	right := round32(newpos.X + pb.Size.X/2)
	top := round32(newpos.Y + pb.Size.Y/2)
	bottom := round32(newpos.Y - pb.Size.Y/2)

	return bodyBounds{
		left, right, top, bottom,
	}
}

var eps float32 = 0.05

func (body *ParticleBody) isCollidingTop(planet worldcollider.WorldCollider, newpos cxmath.Vec2) bool {
	if body.Vel.Y <= 0 {
		return false
	}
	bounds := body.bounds(newpos)

	left := int(round32(body.Pos.X - body.Size.X/2 + eps))
	right := int(round32(body.Pos.X + body.Size.X/2 - eps))
	for x := left; x <= right; x++ {
		if planet.TileIsSolid(x, int(bounds.top)) {
			body.Collisions.Above = true
			return true
		}
	}
	return false
}

func (body *ParticleBody) isCollidingBottom(
	planet worldcollider.WorldCollider,
	newpos cxmath.Vec2,
) bool {
	// don't bother checking if not moving down
	if body.Vel.Y >= 0 {
		return false
	}
	bounds := body.bounds(newpos)
	left := int(round32(body.Pos.X - body.Size.X/2 + eps))
	right := int(round32(body.Pos.X + body.Size.X/2 - eps))
	for x := left; x <= right; x++ {
		colliding := planet.TileIsSolid(x, int(bounds.bottom))
		if colliding {
			body.Collisions.Below = true
			return true
		}
	}
	return false
}

func (body *ParticleBody) isCollidingRight(planet worldcollider.WorldCollider, newpos cxmath.Vec2) bool {
	// don't bother checking if not moving right
	if body.Vel.X <= 0 {
		return false
	}
	bounds := body.bounds(newpos)
	bottom := int(round32(body.Pos.Y - body.Size.Y/2 + eps))
	top := int(round32(body.Pos.Y + body.Size.Y/2 - eps))
	for y := bottom; y <= top; y++ {
		if planet.TileIsSolid(int(bounds.right), y) {
			body.Collisions.Right = true
			return true
		}
	}
	return false
}

func (body *ParticleBody) isCollidingLeft(planet worldcollider.WorldCollider, newpos cxmath.Vec2) bool {
	// don't bother checking if not moving left
	if body.Vel.X >= 0 {
		return false
	}
	bounds := body.bounds(newpos)
	bottom := int(round32(body.Pos.Y - body.Size.Y/2 + eps))
	top := int(round32(body.Pos.Y + body.Size.Y/2 - eps))
	for y := bottom; y <= top; y++ {
		if planet.TileIsSolid(int(bounds.left), y) {
			body.Collisions.Left = true
			return true
		}
	}
	return false
}

func (body *ParticleBody) DetectCollisions(planet worldcollider.WorldCollider, newpos cxmath.Vec2) {
	body.isCollidingTop(planet, newpos)
	body.isCollidingBottom(planet, newpos)
	body.isCollidingRight(planet, newpos)
	body.isCollidingLeft(planet, newpos)
}

func (body *ParticleBody) MoveNoGravity(planet worldcollider.WorldCollider, dt float32) {
	body.PrevPos = body.Pos
	body.Pos = body.Pos.Add(body.Vel.Mult(dt))
}

func (body *ParticleBody) MoveNoBounceGravity(planet worldcollider.WorldCollider, dt float32) {
	body.Collisions.Reset()

	body.PrevPos = body.Pos
	body.PrevVel = body.Vel
	body.Vel.Y += -constants.Gravity * dt

	//todo account drag

	newPos := body.Pos.Add(body.Vel.Mult(dt))

	body.DetectCollisions(planet, newPos)

	if body.Collisions.Above {
		body.Vel.Y = 0
		body.Vel.X = 0
		newPos.Y = body.bounds(newPos).top - 0.5 - body.Size.Y/2
	}
	if body.Collisions.Below {
		body.Vel.Y = 0
		body.Vel.X = 0
		newPos.Y = body.bounds(newPos).bottom + 0.5 + body.Size.Y/2
	}

	if body.Collisions.Left {
		body.Vel.Y = 0
		body.Vel.X = 0
		newPos.X = body.bounds(newPos).left + 0.5 + body.Size.X/2
	}
	if body.Collisions.Right {
		body.Vel.Y = 0
		body.Vel.X = 0
		newPos.X = body.bounds(newPos).right - 0.5 + body.Size.X/2
	}

	body.Pos = newPos
}

func (body *ParticleBody) MoveBounceGravity(planet worldcollider.WorldCollider, dt float32) {
	body.PrevVel = body.Vel
	body.PrevPos = body.Pos
	body.Collisions.Reset()

	body.Vel.Y += -constants.Gravity * dt

	//todo account drag

	newPos := body.Pos.Add(body.Vel.Mult(dt))

	body.DetectCollisions(planet, newPos)

	if body.Collisions.Above {
		body.Vel.Y = 0
		body.Vel.X *= (1 - body.Friction)
		newPos.Y = body.bounds(newPos).top - 0.5 - body.Size.Y/2
	}
	if body.Collisions.Below {
		body.Vel.Y = -body.Vel.Y * body.Elasticity
		body.Vel.X *= (1 - body.Friction)
		newPos.Y = body.bounds(newPos).bottom + 0.5 + body.Size.Y/2
	}
	if body.Collisions.Left {
		body.Vel.X = -body.Vel.X * body.Elasticity
		newPos.X = body.bounds(newPos).left + 0.5 + body.Size.X/2
	}
	if body.Collisions.Right {
		body.Vel.X = -body.Vel.X * body.Elasticity
		newPos.X = body.bounds(newPos).right - 0.5 + body.Size.X/2
	}

	body.Pos = newPos
}

func (body *ParticleBody) MoveNoBounceGravityCallback(planet worldcollider.WorldCollider, dt float32) {
	body.MoveNoBounceGravity(planet, dt)
}
