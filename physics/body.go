package physics

import (
	"github.com/skycoin/cx-game/world"
)

type Body struct {
	Pos  Vec2
	Vel  Vec2
	Size Vec2
}

func (body *Body) Move(planet *world.Planet, dt float32) {
	if body.Vel.IsZero() {
		return
	}

	newPos := body.Pos.Add(body.Vel.Mult(dt))

	if body.Vel.X > 0 { // moving right
		if planet.GetTopLayerTile(int(newPos.X+1.0), int(body.Pos.Y+0.9)).TileType != world.TileTypeNone ||
			planet.GetTopLayerTile(int(newPos.X+1.0), int(body.Pos.Y+0.1)).TileType != world.TileTypeNone {
			newPos.X = float32(int(newPos.X))
			body.Vel.X = 0.0
		}
	} else { // moving left
		if planet.GetTopLayerTile(int(newPos.X), int(body.Pos.Y+0.9)).TileType != world.TileTypeNone ||
			planet.GetTopLayerTile(int(newPos.X), int(body.Pos.Y+0.1)).TileType != world.TileTypeNone {
			newPos.X = float32(int(newPos.X) + 1)
			body.Vel.X = 0.0
		}
	}

	if body.Vel.Y > 0 { // moving up
		if planet.GetTopLayerTile(int(newPos.X), int(newPos.Y+1.0)).TileType != world.TileTypeNone ||
			planet.GetTopLayerTile(int(newPos.X+0.9), int(newPos.Y+1.0)).TileType != world.TileTypeNone {
			newPos.Y = float32(int(newPos.Y))
			body.Vel.Y = 0
		}
	} else { // moving down
		if planet.GetTopLayerTile(int(newPos.X), int(newPos.Y)).TileType != world.TileTypeNone ||
			planet.GetTopLayerTile(int(newPos.X+0.9), int(newPos.Y)).TileType != world.TileTypeNone {
			newPos.Y = float32(int(newPos.Y) + 1)
			body.Vel.Y = 0
		}
	}

	body.Pos = newPos
}
