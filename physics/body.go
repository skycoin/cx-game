package physics

import (
	"math"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/world"
)

// epsilon parameter for values that are "close enough"
const eps = 0.05

type Body struct {
	Pos            Vec2
	Vel            Vec2
	Size           Vec2
	collidingLines []float32
}

type bodyBounds struct {
	left, right, top, bottom                 float32
	leftTile, rightTile, topTile, bottomTile int
}

func (body Body) bounds(newpos Vec2) bodyBounds {
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

func (body *Body) isCollidingLeft(planet *world.Planet, newpos Vec2) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving left
	if body.Vel.X >= 0 {
		return false
	}
	bottom := int(round32(body.Pos.Y - body.Size.Y/2 + eps))
	top := int(round32(body.Pos.Y + body.Size.Y/2 - eps))
	for y := bottom; y <= top; y++ {
		tile := planet.GetTopLayerTile(bounds.leftTile, y)
		if tile.TileType != world.TileTypeNone {
			return true
		}
	}
	return false
}

func (body *Body) isCollidingRight(planet *world.Planet, newpos Vec2) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving right
	if body.Vel.X <= 0 {
		return false
	}
	bottom := int(round32(body.Pos.Y - body.Size.Y/2 + eps))
	top := int(round32(body.Pos.Y + body.Size.Y/2 - eps))
	for y := bottom; y <= top; y++ {
		tile := planet.GetTopLayerTile(bounds.rightTile, y)
		if tile.TileType != world.TileTypeNone {
			return true
		}
	}
	return false
}

func (body *Body) isCollidingTop(planet *world.Planet, newpos Vec2) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving up
	if body.Vel.Y <= 0 {
		return false
	}
	left := int(round32(body.Pos.X - body.Size.X/2 + eps))
	right := int(round32(body.Pos.X + body.Size.X/2 - eps))
	for x := left; x <= right; x++ {
		tile := planet.GetTopLayerTile(x, bounds.topTile)
		if tile.TileType != world.TileTypeNone {
			return true
		}
	}
	return false
}

func (body *Body) isCollidingBottom(planet *world.Planet, newpos Vec2) bool {
	bounds := body.bounds(newpos)
	// don't bother checking if not moving down
	if body.Vel.Y >= 0 {
		return false
	}
	left := int(round32(body.Pos.X - body.Size.X/2 + eps))
	right := int(round32(body.Pos.X + body.Size.X/2 - eps))
	for x := left; x <= right; x++ {
		tile := planet.GetTopLayerTile(x, bounds.bottomTile)
		if tile.TileType != world.TileTypeNone {
			return true
		}
	}
	return false
}

func (body *Body) Move(planet *world.Planet, dt float32) {
	body.collidingLines = []float32{}

	if body.Vel.IsZero() {
		return
	}

	newPos := body.Pos.Add(body.Vel.Mult(dt))

	if body.isCollidingLeft(planet, newPos) {
		body.Vel.X = 0
		newPos.X = float32(body.bounds(newPos).leftTile) + 0.5 + body.Size.X/2

		body.collidingLines = append(body.collidingLines, []float32{
			newPos.X - body.Size.X/2, newPos.Y - body.Size.Y/2, 0.0,
			newPos.X - body.Size.X/2, newPos.Y + body.Size.Y/2, 0.0,
		}...)
	}
	if body.isCollidingRight(planet, newPos) {
		body.Vel.X = 0
		newPos.X = float32(body.bounds(newPos).rightTile) - 0.5 - body.Size.X/2

		body.collidingLines = append(body.collidingLines, []float32{
			newPos.X + body.Size.X/2, newPos.Y - body.Size.Y/2, 0.0,
			newPos.X + body.Size.X/2, newPos.Y + body.Size.Y/2, 0.0,
		}...)
	}
	if body.isCollidingTop(planet, newPos) {
		body.Vel.Y = 0
		newPos.Y = float32(body.bounds(newPos).topTile) - 0.5 - body.Size.Y/2

		body.collidingLines = append(body.collidingLines, []float32{
			newPos.X - body.Size.X/2, newPos.Y + body.Size.Y/2, 0.0,
			newPos.X + body.Size.X/2, newPos.Y + body.Size.Y/2, 0.0,
		}...)
	}
	if body.isCollidingBottom(planet, newPos) {
		body.Vel.Y = 0
		newPos.Y = float32(body.bounds(newPos).bottomTile) + 0.5 + body.Size.Y/2

		body.collidingLines = append(body.collidingLines, []float32{
			newPos.X - body.Size.X/2, newPos.Y - body.Size.Y/2, 0.0,
			newPos.X + body.Size.X/2, newPos.Y - body.Size.Y/2, 0.0,
		}...)
	}

	body.Pos = newPos
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

func (body *Body) GetCollidingLines(cam *camera.Camera) []float32 {
	collidingLines := []float32{}

	for i := 0; i < len(body.collidingLines); i += 3 {
		collidingLines = append(collidingLines, []float32{
			body.collidingLines[i] - cam.X,
			body.collidingLines[i+1] - cam.Y,
			0.0,
		}...)
	}

	return collidingLines
}
