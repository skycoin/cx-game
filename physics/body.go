package physics

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world/worldcollider"
)

// epsilon parameter for values that are "close enough"
const eps = 0.05

type DamageFunc func(damage int)
type Body struct {
	Pos  cxmath.Vec2
	Vel  cxmath.Vec2
	Size cxmath.Vec2

	PreviousTransform     mgl32.Mat4
	InterpolatedTransform mgl32.Mat4

	Collisions     CollisionInfo
	collidingLines []float32

	Damage  DamageFunc
	Deleted bool
}

func (body Body) Transform() mgl32.Mat4 {
	return mgl32.Translate3D(body.Pos.X, body.Pos.Y, 0)
}

func (body *Body) SavePreviousTransform() {
	body.PreviousTransform = body.Transform()
}

func (body *Body) UpdateInterpolatedTransform(alpha float32) {
	prevPart := body.PreviousTransform.Mul(1 - alpha)
	nextPart := body.Transform().Mul(alpha)
	body.InterpolatedTransform = prevPart.Add(nextPart)
}

type bodyBounds struct {
	left, right, top, bottom                 float32
	leftTile, rightTile, topTile, bottomTile int
}

func (body Body) bounds(newpos cxmath.Vec2) bodyBounds {
	left := round32(newpos.X - body.Size.X/2)
	leftTile := int(left)
	right := round32(newpos.X + body.Size.X/2)
	rightTile := int(right)
	top := round32(newpos.Y + body.Size.Y/2)
	topTile := int(top)
	bottom := round32(newpos.Y - body.Size.Y/2)
	bottomTile := int(bottom)
	return bodyBounds{
		left: left, right: right, top: top, bottom: bottom,
		leftTile: leftTile, rightTile: rightTile,
		topTile: topTile, bottomTile: bottomTile,
	}
}

func round32(x float32) float32 {
	return float32(math.Round(float64(x)))
}

func (body *Body) isCollidingLeft(collider worldcollider.WorldCollider, newpos cxmath.Vec2) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving left
	if body.Vel.X >= 0 {
		return false
	}
	bottom := int(round32(body.Pos.Y - body.Size.Y/2 + eps))
	top := int(round32(body.Pos.Y + body.Size.Y/2 - eps))
	for y := bottom; y <= top; y++ {
		if collider.TileIsSolid(bounds.leftTile, y) {
			return true
		}
	}
	return false
}

func (body *Body) isCollidingRight(collider worldcollider.WorldCollider, newpos cxmath.Vec2) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving right
	if body.Vel.X <= 0 {
		return false
	}
	bottom := int(round32(body.Pos.Y - body.Size.Y/2 + eps))
	top := int(round32(body.Pos.Y + body.Size.Y/2 - eps))
	for y := bottom; y <= top; y++ {
		if collider.TileIsSolid(bounds.rightTile, y) {
			return true
		}
	}
	return false
}

func (body *Body) isCollidingTop(collider worldcollider.WorldCollider, newpos cxmath.Vec2) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving up
	if body.Vel.Y <= 0 {
		return false
	}
	left := int(round32(body.Pos.X - body.Size.X/2 + eps))
	right := int(round32(body.Pos.X + body.Size.X/2 - eps))
	for x := left; x <= right; x++ {
		if collider.TileIsSolid(x, bounds.topTile) {
			return true
		}
	}
	return false
}

func (body *Body) isCollidingBottom(collider worldcollider.WorldCollider, newpos cxmath.Vec2) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving down
	if body.Vel.Y >= 0 {
		return false
	}
	left := int(round32(body.Pos.X - body.Size.X/2 + eps))
	right := int(round32(body.Pos.X + body.Size.X/2 - eps))
	for x := left; x <= right; x++ {
		if collider.TileIsSolid(x, bounds.bottomTile) {
			return true
		}
	}
	return false
}

func (body *Body) Move(collider worldcollider.WorldCollider, dt float32) {
	body.collidingLines = []float32{}
	body.Collisions.Reset()

	if body.Vel.IsZero() {
		return
	}

	newPos := body.Pos.Add(body.Vel.Mult(dt))

	if body.isCollidingLeft(collider, newPos) {
		body.Collisions.Left = true
		body.Vel.X = 0
		newPos.X = float32(body.bounds(newPos).leftTile) + 0.5 + body.Size.X/2

		body.collidingLines = append(body.collidingLines, []float32{
			newPos.X - body.Size.X/2, newPos.Y - body.Size.Y/2, 0.0,
			newPos.X - body.Size.X/2, newPos.Y + body.Size.Y/2, 0.0,
		}...)
	}
	if body.isCollidingRight(collider, newPos) {
		body.Collisions.Right = true
		body.Vel.X = 0
		newPos.X = float32(body.bounds(newPos).rightTile) - 0.5 - body.Size.X/2

		body.collidingLines = append(body.collidingLines, []float32{
			newPos.X + body.Size.X/2, newPos.Y - body.Size.Y/2, 0.0,
			newPos.X + body.Size.X/2, newPos.Y + body.Size.Y/2, 0.0,
		}...)
	}
	if body.isCollidingTop(collider, newPos) {
		body.Collisions.Above = true
		body.Vel.Y = 0
		newPos.Y = float32(body.bounds(newPos).topTile) - 0.5 - body.Size.Y/2

		body.collidingLines = append(body.collidingLines, []float32{
			newPos.X - body.Size.X/2, newPos.Y + body.Size.Y/2, 0.0,
			newPos.X + body.Size.X/2, newPos.Y + body.Size.Y/2, 0.0,
		}...)
	}
	if body.isCollidingBottom(collider, newPos) {
		body.Collisions.Below = true
		body.Vel.Y = 0
		newPos.Y = float32(body.bounds(newPos).bottomTile) + 0.5 + body.Size.Y/2

		body.collidingLines = append(body.collidingLines, []float32{
			newPos.X - body.Size.X/2, newPos.Y - body.Size.Y/2, 0.0,
			newPos.X + body.Size.X/2, newPos.Y - body.Size.Y/2, 0.0,
		}...)
	}

	newPosMgl32 := mgl32.Vec2{newPos.X, newPos.Y}
	offset := collider.WrapAroundOffset(newPosMgl32)
	newPosMgl32 = newPosMgl32.Add(offset)
	body.Pos = cxmath.Vec2{newPosMgl32.X(), newPosMgl32.Y()}
	// move previous transform to avoid weird interpolation around boundaries
	body.PreviousTransform = body.PreviousTransform.
		Mul4(mgl32.Translate3D(offset.X(), offset.Y(), 0))
}

func (body *Body) GetBBoxLines() []float32 {
	x := body.Pos.X - body.Size.X/2
	y := body.Pos.Y - body.Size.Y/2
	return []float32{
		x, y, 0.0,
		x + body.Size.X, y, 0.0,

		x + body.Size.X, y, 0.0,
		x + body.Size.X, y + body.Size.Y, 0.0,

		x + body.Size.X, y + body.Size.Y, 0.0,
		x, y + body.Size.Y, 0.0,

		x, y + body.Size.Y, 0.0,
		x, y, 0.0,
	}
}

func (body *Body) GetCollidingLines() []float32 {
	collidingLines := []float32{}

	for i := 0; i < len(body.collidingLines); i += 3 {
		collidingLines = append(collidingLines, []float32{
			body.collidingLines[i],
			body.collidingLines[i+1],
			0.0,
		}...)
	}

	return collidingLines
}
