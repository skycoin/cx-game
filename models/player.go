package models

import (
	"image"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
)

type Player struct {
	physics.Body
	RGBA      *image.RGBA
	ImageSize image.Point
	movSpeed  float32
	jumpSpeed float32
	spriteId  int
}

func NewPlayer() *Player {
	spriteId := spriteloader.LoadSingleSprite(
		"./assets/character/character.png", "player" )
	player := Player{
		Body: physics.Body{
			Size: physics.Vec2{X: 2.0 * 64/96, Y: 2.0},
		},
		movSpeed:  3.0,
		jumpSpeed: 12.0,
		spriteId:  spriteId,
	}
	return &player
}

func (player *Player) Draw(cam *camera.Camera) {

	x := player.Pos.X - cam.X
	y := player.Pos.Y - cam.Y
	if !cam.IsInBoundsF(player.Pos.X, player.Pos.Y) {
		return
	}

	spriteloader.DrawSpriteQuad(
		x, y,
		player.Size.X, player.Size.Y, player.spriteId,
	)
}

func boolToFloat(x bool) float32 {
	if x {
		return 1
	} else {
		return 0
	}
}

func (player *Player) Tick(leftPressed, rightPressed, spacePressed bool, planet *world.Planet, dt float32) {
	player.Vel.X = (boolToFloat(rightPressed) - boolToFloat(leftPressed)) * player.movSpeed
	player.Vel.Y -= physics.Gravity * dt

	if spacePressed {
		player.Vel.Y = player.jumpSpeed
	}

	player.Move(planet, dt)
}
