package physics

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world/worldcollider"
)

// epsilon parameter for values that are "close enough"
const eps = 0.05

type DamageFunc func(damage int)
type Body struct {
	Pos       cxmath.Vec2
	Vel       cxmath.Vec2
	Size      cxmath.Vec2
	Direction float32

	PreviousTransform     mgl32.Mat4
	InterpolatedTransform mgl32.Mat4

	Collisions CollisionInfo

	Damage  DamageFunc
	Deleted bool

	IsIgnoringPlatforms bool
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

func (body *Body) isCollidingBottom(
	collider worldcollider.WorldCollider,
	newpos cxmath.Vec2,
) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving down
	if body.Vel.Y >= 0 {
		return false
	}
	left := int(round32(body.Pos.X - body.Size.X/2 + eps))
	right := int(round32(body.Pos.X + body.Size.X/2 - eps))
	for x := left; x <= right; x++ {
		colliding := collider.
			TileTopIsSolid(x, bounds.bottomTile, body.IsIgnoringPlatforms)
		if colliding {
			return true
		}
	}
	return false
}

func (body *Body) Move(collider worldcollider.WorldCollider, dt float32) {
	body.Collisions.Reset()

	body.Vel.Y -= Gravity * dt

	newPos := body.Pos.Add(body.Vel.Mult(dt))

	if body.isCollidingLeft(collider, newPos) {
		body.Collisions.Left = true
		body.Vel.X = 0
		newPos.X = body.bounds(newPos).left + 0.5 + body.Size.X/2

	}
	if body.isCollidingRight(collider, newPos) {
		body.Collisions.Right = true
		body.Vel.X = 0
		newPos.X = body.bounds(newPos).right - 0.5 - body.Size.X/2

	}
	if body.isCollidingTop(collider, newPos) {
		body.Collisions.Above = true
		body.Vel.Y = 0
		newPos.Y = body.bounds(newPos).top - 0.5 - body.Size.Y/2
	}
	if body.isCollidingBottom(collider, newPos) {
		body.Collisions.Below = true
		body.Vel.Y = 0
		newPos.Y = body.bounds(newPos).bottom + 0.5 + body.Size.Y/2
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
		//bottom
		x, y,
		x + body.Size.X, y,

		//right
		x + body.Size.X, y,
		x + body.Size.X, y + body.Size.Y,

		//top
		x + body.Size.X, y + body.Size.Y,
		x, y + body.Size.Y,

		//left
		x, y + body.Size.Y,
		x, y,
	}
}

func (body *Body) GetCollidingLines() []float32 {
	collidingLines := make([]float32, 0, 16)
	bboxLines := body.GetBBoxLines()
	if body.Collisions.Below {
		collidingLines = append(collidingLines, bboxLines[0:4]...)
	}
	if body.Collisions.Right {
		collidingLines = append(collidingLines, bboxLines[4:8]...)
	}
	if body.Collisions.Above {
		collidingLines = append(collidingLines, bboxLines[8:12]...)
	}
	if body.Collisions.Left {
		collidingLines = append(collidingLines, bboxLines[12:16]...)
	}
	fmt.Println(collidingLines)

	// for i := 0; i < len(body.collidingLines); i += 2 {
	// 	collidingLines = append(collidingLines, []float32{
	// 		body.collidingLines[i],
	// 		body.collidingLines[i+1],
	// 	}...)
	// }

	return collidingLines
}

func (body *Body) IsOnGround() bool {
	return body.Collisions.Below
}
