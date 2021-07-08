package enemies

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/enemies/pathfinding"
)

type Enemy struct {
	physics.Body
	SpriteID uint32
	Health int
	TimeSinceLastJump float32
	PathfindingBehaviourID pathfinding.BehaviourID
}

func InitBasicEnemies() {
	basicEnemySpriteId = spriteloader.
		LoadSingleSprite("./assets/enemies/basic-enemy.png", "basic-enemy")
}

// TODO load an actual sprite here
var basicEnemySpriteId spriteloader.SpriteID
var basicEnemyMovSpeed = float32(1)
var basicEnemies = []*Enemy{}

const enemyJumpVelocity = 15

// TODO create a system to handle projectiles, melee attacks, etc
var playerStrikeRange = float32(1)

func TickBasicEnemies(
	world *world.Planet, dt float32,
	player *models.Player, playerIsAttacking bool,
) {
	nextEnemies := []*Enemy{}
	for idx, _ := range basicEnemies {
		enemy := basicEnemies[idx]
		enemy.Tick(world, dt, player, playerIsAttacking)
		if enemy.Health > 0 {
			nextEnemies = append(nextEnemies, enemy)
		} else {
			enemy.Deleted = true
		}
	}
	basicEnemies = nextEnemies
}

func DrawBasicEnemies(cam *camera.Camera) {
	for _, enemy := range basicEnemies {
		if cam.IsInBounds(int(enemy.Pos.X), int(enemy.Pos.Y)) {
			enemy.Draw(cam)
		}
	}
}

func SpawnBasicEnemy(x, y float32) {
	enemy := Enemy{
		Body: physics.Body{
			Size: cxmath.Vec2{X: 3.0, Y: 3.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		Health: 5,
		SpriteID: uint32(basicEnemySpriteId),
		PathfindingBehaviourID: pathfinding.WalkingBehaviourID,
	}
	enemy.Damage = func(damage int) {
		enemy.Health -= 1
	}
	physics.RegisterBody(&enemy.Body)
	basicEnemies = append(basicEnemies, &enemy)
}

func SpawnLeapingEnemy(x,y float32) {
	enemy := Enemy{
		Body: physics.Body{
			Size: cxmath.Vec2{X:2.0, Y: 2.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		Health:5,
		// TODO swap out sprite
		SpriteID: uint32(basicEnemySpriteId),
		PathfindingBehaviourID: pathfinding.LeapingBehaviourID,
	}
	physics.RegisterBody(&enemy.Body)
	basicEnemies = append(basicEnemies, &enemy)
}

func (enemy Enemy) isStuck() bool {
	return (enemy.Collisions.Left || enemy.Collisions.Right) &&
		!enemy.Collisions.Below
}

func (enemy *Enemy) Tick(
	world *world.Planet, dt float32, player *models.Player,
	playerIsAttacking bool,
) bool {
	enemy.PathfindingBehaviourID.Get().Follow(pathfinding.BehaviourContext{
		Self: &enemy.Body,
		PlayerPos: player.Pos.Mgl32(),
	})

	playerIsCloseEnoughToStrike :=
		player.Pos.Sub(enemy.Pos).LengthSqr() <
			playerStrikeRange*playerStrikeRange

	stillAlive := !playerIsAttacking || !playerIsCloseEnoughToStrike
	return stillAlive
}

func (enemy *Enemy) Draw(cam *camera.Camera) {
	camX := enemy.Pos.X - cam.X
	camY := enemy.Pos.Y - cam.Y

	spriteloader.DrawSpriteQuad(
		camX, camY,
		enemy.Size.X, enemy.Size.Y,
		basicEnemySpriteId,
	)
}
