package physics

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/world"
)

type Body struct {
	Pos            Vec2
	Vel            Vec2
	Size           Vec2
	collidingLines []float32
}

func (body *Body) Move(planet *world.Planet, dt float32) {
	body.collidingLines = []float32{}

	if body.Vel.IsZero() {
		return
	}

	newPos := body.Pos.Add(body.Vel.Mult(dt))

	if body.Vel.X > 0 { // moving right
		if planet.GetTopLayerTile(int(newPos.X+1.0), int(body.Pos.Y+0.9)).TileType != world.TileTypeNone ||
			planet.GetTopLayerTile(int(newPos.X+1.0), int(body.Pos.Y)).TileType != world.TileTypeNone {
			newPos.X = float32(int(newPos.X))
			body.Vel.X = 0.0

			body.collidingLines = append(body.collidingLines, []float32{
				newPos.X + 1.0, newPos.Y, 0.0,
				newPos.X + 1.0, newPos.Y + 1.0, 0.0,
			}...)
		}
	} else if body.Vel.X < 0 { // moving left
		if planet.GetTopLayerTile(int(newPos.X), int(body.Pos.Y+0.9)).TileType != world.TileTypeNone ||
			planet.GetTopLayerTile(int(newPos.X), int(body.Pos.Y)).TileType != world.TileTypeNone {
			newPos.X = float32(int(newPos.X) + 1)
			body.Vel.X = 0.0

			body.collidingLines = append(body.collidingLines, []float32{
				newPos.X, newPos.Y, 0.0,
				newPos.X, newPos.Y + 1.0, 0.0,
			}...)
		}
	}

	if body.Vel.Y > 0 { // moving up
		if planet.GetTopLayerTile(int(newPos.X+0.05), int(newPos.Y+1.0)).TileType != world.TileTypeNone ||
			planet.GetTopLayerTile(int(newPos.X+0.95), int(newPos.Y+1.0)).TileType != world.TileTypeNone {
			newPos.Y = float32(int(newPos.Y))
			body.Vel.Y = 0

			body.collidingLines = append(body.collidingLines, []float32{
				newPos.X, newPos.Y + 1.0, 0.0,
				newPos.X + 1.0, newPos.Y + 1.0, 0.0,
			}...)

		}
	} else if body.Vel.Y < 0 { // moving down
		if planet.GetTopLayerTile(int(newPos.X+0.05), int(newPos.Y)).TileType != world.TileTypeNone ||
			planet.GetTopLayerTile(int(newPos.X+0.95), int(newPos.Y)).TileType != world.TileTypeNone {
			newPos.Y = float32(int(newPos.Y) + 1)
			body.Vel.Y = 0

			body.collidingLines = append(body.collidingLines, []float32{
				newPos.X, newPos.Y, 0.0,
				newPos.X + 1.0, newPos.Y, 0.0,
			}...)

			// check the sides and correct the position to be centered when fall near another tile
			// this is because we are checking with a slightly smaller width
			if planet.GetTopLayerTile(int(newPos.X+1.0), int(body.Pos.Y+0.9)).TileType != world.TileTypeNone ||
				planet.GetTopLayerTile(int(newPos.X+1.0), int(body.Pos.Y)).TileType != world.TileTypeNone {
				newPos.X = float32(int(newPos.X))
				body.Vel.X = 0.0
			}
			if planet.GetTopLayerTile(int(newPos.X), int(body.Pos.Y+0.9)).TileType != world.TileTypeNone ||
				planet.GetTopLayerTile(int(newPos.X), int(body.Pos.Y)).TileType != world.TileTypeNone {
				newPos.X = float32(int(newPos.X) + 1)
				body.Vel.X = 0.0
			}
		}
	}

	body.Pos = newPos
}

func (body *Body) GetBBoxLinesRelative(cam *camera.Camera) []float32 {
	x := body.Pos.X - cam.X - body.Size.X/2
	y := body.Pos.Y - cam.Y - body.Size.Y/2
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

func (body *Body) GetCollidingLinesRelative(cam *camera.Camera) []float32 {
	collidingLines := []float32{}

	for i := 0; i < len(body.collidingLines); i += 3 {
		collidingLines = append(collidingLines, []float32{
			body.collidingLines[i] - cam.X - 0.5,
			body.collidingLines[i+1] - cam.Y - 0.5,
			0.0,
		}...)
	}

	return collidingLines
}
