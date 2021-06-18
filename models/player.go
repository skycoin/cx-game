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
	physics.RegisterBody(&player.Body)

	return &player
}

func (player *Player) Draw(cam *camera.Camera, planet *world.Planet) {

	worldTransform := player.InterpolatedTransform
	worldPos := worldTransform.Col(3).Vec2()

	disp := planet.ShortestDisplacement(
		mgl32.Vec2{cam.X, cam.Y},
		worldPos )

	spriteloader.DrawSpriteQuad(
		disp.X(), disp.Y(),
		player.Size.X, player.Size.Y, player.spriteId,
	)
}

func (player *Player) FixedTick(controlled bool, planet *world.Planet) {
	player.Vel.Y -= physics.Gravity * physics.TimeStep

	if controlled {
		player.Vel.X = input.GetAxis(input.HORIZONTAL) * player.movSpeed
	} else {
		player.Vel.X = 0
	}
}

func (player *Player) Jump() (didJump bool) {
	if player.Vel.Y != 0 {
		return false
	}
	player.Vel.Y += player.jumpSpeed
	return true
}
