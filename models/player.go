package models

import (
	"image"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/input"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/world"
)

type Player struct {
	physics.Body
	meta            agents.AgentMeta
	RGBA            *image.RGBA
	ImageSize       image.Point
	acceleration    float32
	AdditionalJumps uint
	spriteId        int
	MovementType    MovementType
	XDirection float32 // 1 when facing right, -1 when facing left
}

type MovementType uint

const (
	NORMAL MovementType = iota
	WALL_SLIDING
	FLYING
)

var (
	maxVerticalSpeed           float32 = 5
	minJumpSpeed, maxJumpSpeed float32
	maxAdditionalJumps         uint = 1
	jumpCounter                uint
	wallSlideSpeed             float32 = -5
)

func NewPlayer() *Player {
	spriteId := spriteloader.LoadSingleSprite(
		"./assets/character/character.png", "player")
	player := Player{
		Body: physics.Body{
			Size: physics.Vec2{X: 2, Y: 3},
		},
		acceleration: 2.0,
		spriteId:     spriteId,
		meta: agents.AgentMeta{
			MovementSpeed:   7,
			MinJumpHeight:   1,
			MaxJumpHeight:   7,
			DynamicFriction: 0.55,
		},
		XDirection: 1, // start facing right
	}
	player.meta.DynamicFriction = utility.ClampF(player.meta.DynamicFriction, 0, 1)
	minJumpSpeed = cxmath.Sqrt(2 * cxmath.Abs(physics.Gravity) * player.meta.MinJumpHeight)
	maxJumpSpeed = cxmath.Sqrt(2 * cxmath.Abs(physics.Gravity) * player.meta.MaxJumpHeight)
	physics.RegisterBody(&player.Body)

	return &player
}

func (player *Player) Draw(cam *camera.Camera, planet *world.Planet) {

	worldTransform := player.InterpolatedTransform
	worldPos := worldTransform.Col(3).Vec2()

	disp := planet.ShortestDisplacement(
		mgl32.Vec2{cam.X, cam.Y},
		worldPos)

	spriteloader.DrawSpriteQuad(
		disp.X(), disp.Y(),
		// player sprite actually faces left so throw an extra (-) here
		-player.Size.X * player.XDirection, player.Size.Y, player.spriteId,
	)
}

func (player *Player) FixedTick(controlled bool, planet *world.Planet) {
	if player.Collisions.Below {
		player.MovementType = NORMAL
	} else if player.Collisions.Left || player.Collisions.Right {
		player.MovementType = WALL_SLIDING
	} else {
		if player.MovementType == WALL_SLIDING {
			player.MovementType = NORMAL
		}
	}
	if !player.Collisions.Below && input.GetButtonUp("jump") {
		if player.Vel.Y > minJumpSpeed {
			player.Vel.Y = minJumpSpeed
		}
	}
	player.Vel.Y -= physics.Gravity * physics.TimeStep

	if controlled {
		inputXAxis := input.GetAxis(input.HORIZONTAL)
		player.Vel.X += inputXAxis * player.acceleration
		if inputXAxis != 0 {
			player.XDirection = math32.Sign(inputXAxis)
		}
	}
	if player.Vel.X != 0 {
		friction := cxmath.Sign(player.Vel.X) * -1 * player.acceleration * player.meta.DynamicFriction
		if cxmath.Abs(player.Vel.X) <= player.acceleration*player.meta.DynamicFriction && input.GetAxis(input.HORIZONTAL) == 0 {
			player.Vel.X = 0
		} else {
			player.Vel.X += friction
		}

	}
	player.Vel.X = utility.ClampF(player.Vel.X, -player.meta.MovementSpeed, player.meta.MovementSpeed)
	player.ApplyMovementConstraints(controlled)
}

func (player *Player) Jump() bool {
	if !player.Collisions.Below && (player.Collisions.Left || player.Collisions.Right) {
		player.Vel.Y = maxJumpSpeed
		// player.Vel.X = cxmath.Min(-player.Vel.X, 15*cxmath.Sign(player.Vel.X))
		player.Vel.X = cxmath.Sign(input.GetAxis(input.HORIZONTAL)) * -1 * 15
		// // fmt.Println(cxmath.Sign(player.Vel.X))
		// jumpCounter = maxAdditionalJumps
		return true
	}
	if player.Collisions.Below {
		// fmt.Println(jumpSpeed)
		player.Vel.Y = maxJumpSpeed
		jumpCounter = maxAdditionalJumps
		return true
	}

	if jumpCounter > 0 {
		jumpCounter -= 1
		player.Vel.Y = maxJumpSpeed
		return true
	}

	return false
}

func (player *Player) ToggleFlying() {
	if player.MovementType == FLYING {
		player.MovementType = NORMAL
	} else {
		player.MovementType = FLYING
	}
}

//states - running
func (player *Player) ApplyMovementConstraints(controlled bool) {
	switch player.MovementType {
	case NORMAL: // moving | idle
		//sprite normal
	case WALL_SLIDING: // wall slide
		player.Vel.Y = cxmath.Max(player.Vel.Y, -6)
		player.AdditionalJumps = maxAdditionalJumps
	case FLYING: // flying

		if input.GetAxis(input.VERTICAL) == 0 || !controlled {
			player.Vel.Y = -3
		} else {
			player.Vel.Y = input.GetAxis(input.VERTICAL) * maxVerticalSpeed
			// player.Vel.Y = utility.ClampF(player.Vel.Y, -maxVerticalSpeed, maxVerticalSpeed)
		}
	}
}

func (m MovementType) String() string {
	switch m {
	case NORMAL:
		return "moving"
	case WALL_SLIDING:
		return "wall sliding"
	case FLYING:
		return "flying"
	default:
		return "unknown"
	}
}
