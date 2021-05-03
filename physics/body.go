package physics

import (
	"log"

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
		body.Pos.X = newPos.X
	} else {
		body.Pos.X = newPos.X
	}

	if body.Vel.Y > 0 { // moving up
		body.Pos.Y = newPos.Y
	} else { // moving down
		if planet.GetTopTile(int(newPos.X), int(newPos.Y-1.0)).TileType != world.TileTypeNone ||
			planet.GetTopTile(int(newPos.X+0.9), int(newPos.Y+1.0)).TileType != world.TileTypeNone {
			newPos.Y = float32(int(newPos.Y))
			body.Vel.Y = 0
			log.Println("on ground")
		}
	}

	//topRight := body.Pos.Add(body.Size.Mult(0.4999))
	//bottomLeft := body.Pos.Add(body.Size.Mult(-0.4999))
}
