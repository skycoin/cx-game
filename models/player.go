package models

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/input"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/physics/movement"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/world"
)

type Player struct {
	physics.Body
	movement.MovementComponent
	Controlled bool
	// RGBA            *image.RGBA
	// ImageSize       image.Point
	helmId        int
	suitId        int
	helmSpriteIds [4]int
	suitSpriteIds [4]int
	XDirection    float32 // 1 when facing right, -1 when facing left
}

func NewPlayer() *Player {
	// spriteId := spriteloader.LoadSingleSprite(
	// 	"./assets/character/character.png", "player")
	helmSpriteSheetId = spriteloader.LoadSpriteSheet("./assets/character/character-helmets.png")
	suitSpriteSheetId = spriteloader.LoadSpriteSheet("./assets/character/character-suits.png")

	player := Player{
		Body: physics.Body{
			Size: physics.Vec2{X: 2, Y: 3},
		},
		MovementComponent: movement.NewPlayerMovementComponent(),

		XDirection: 1, // start facing right
		Controlled: true,
	}

	if !outfitsLoaded {
		outfitsLoaded = true
		loadOutfits()
	}
	player.SetHelm(DEFAULT_HELM)
	player.SetSuit(DEFAULT_SUIT)

	maxJumpSpeed = cxmath.Sqrt(2 * cxmath.Abs(physics.Gravity) * player.MovementMeta.Jumpheight)
	minJumpSpeed = maxJumpSpeed / 4
	physics.RegisterBody(&player.Body)

	return &player
}

func (player *Player) Draw(cam *camera.Camera, planet *world.Planet) {
	disp := planet.ShortestDisplacement(
		mgl32.Vec2{cam.X, cam.Y},
		player.InterpolatedTransform.Col(3).Vec2())

	player.DrawOutfit(disp)

}

func (player *Player) FixedTick(planet *world.Planet) {
	//todo separate more logic
	//
	player.MovementBeforeTick()

	if player.Controlled {
		inputXAxis := input.GetAxis(input.HORIZONTAL)
		player.Vel.X += inputXAxis * player.MovementMeta.Acceleration * player.ActiveMovementType.GetMovementSpeedModifier()
		if inputXAxis != 0 {
			player.XDirection = math32.Sign(inputXAxis)
		}

		if player.ActiveMovementType == movement.FLYING {
			inputYAxis := input.GetAxis(input.VERTICAL)
			player.Vel.Y = inputYAxis * maxVerticalSpeed
		}
	}
	player.Vel.Y -= physics.Gravity * physics.TimeStep * player.ActiveMovementType.GetGravityModifier()

	if player.Vel.X != 0 {
		friction := cxmath.Sign(player.Vel.X) * player.MovementMeta.Acceleration * player.MovementMeta.DynamicFriction * player.ActiveMovementType.GetFrictionModifier()

		//to stop player from jiggling
		if cxmath.Abs(player.Vel.X) <= player.MovementMeta.Acceleration*player.MovementMeta.DynamicFriction*player.ActiveMovementType.GetFrictionModifier() && input.GetAxis(input.HORIZONTAL) == 0 {
			player.Vel.X = 0
		} else {
			player.Vel.X -= friction
			// fmt.Println(player.Vel.X)
		}
	}
	player.Vel.X = utility.ClampF(player.Vel.X, -player.MovementMeta.MovSpeed*player.ActiveMovementType.GetMovementSpeedModifier(), player.MovementMeta.MovSpeed*player.ActiveMovementType.GetMovementSpeedModifier())
	player.MovementAfterTick(planet)
}
