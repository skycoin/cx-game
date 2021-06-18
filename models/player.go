package models

import (
	"image"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/input"
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
		"./assets/character/character.png", "player")
	player := Player{
		Body: physics.Body{
			Size: physics.Vec2{X: 2.0 * 64 / 96, Y: 2.0},
		},
		movSpeed:  3.0,
		jumpSpeed: 12.0,
		spriteId:  spriteId,
	}
	return &player
}

func (player *Player) Draw(cam *camera.Camera, planet *world.Planet) {

	disp := planet.ShortestDisplacement(
		mgl32.Vec2{cam.X, cam.Y},
		mgl32.Vec2{player.Pos.X, player.Pos.Y})

	spriteloader.DrawSpriteQuad(
		disp.X(), disp.Y(),
		player.Size.X, player.Size.Y, player.spriteId,
	)
}

func (player *Player) Tick(controlled bool, planet *world.Planet, dt float32) {
	player.Vel.Y -= physics.Gravity * dt

	if controlled {
		player.Vel.X = input.GetAxis(input.HORIZONTAL) * player.movSpeed
	} else {
		player.Vel.X = 0
	}

	player.Move(planet, dt)
}

func (player *Player) Jump() bool {
	if player.Vel.Y != 0 {
		return false
	}
	player.Vel.Y += player.jumpSpeed
	return true
}
