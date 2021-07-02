package enemies

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
)

type BasicEnemy struct {
	physics.Body
	Health int
}

func InitBasicEnemies() {
	basicEnemySpriteId = spriteloader.
		LoadSingleSprite("./assets/enemies/basic-enemy.png", "basic-enemy")
}

// TODO load an actual sprite here
var basicEnemySpriteId int
var basicEnemyMovSpeed = float32(1)
var basicEnemies = []*BasicEnemy{}

// TODO create a system to handle projectiles, melee attacks, etc
var playerStrikeRange = float32(1)

func TickBasicEnemies(
	world *world.Planet, dt float32,
	player *models.Player, playerIsAttacking bool,
) {
	nextEnemies := []*BasicEnemy{}
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
	enemy := BasicEnemy{
		Body: physics.Body{
			Size: cxmath.Vec2{X: 3.0, Y: 3.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		Health: 5,
	}
	enemy.Damage = func(damage int) {
		enemy.Health -= 1
	}
	physics.RegisterBody(&enemy.Body)
	basicEnemies = append(basicEnemies, &enemy)
}

func sign(x float32) float32 {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

func (enemy *BasicEnemy) Tick(
	world *world.Planet, dt float32, player *models.Player,
	playerIsAttacking bool,
) bool {
	enemy.Vel.X = basicEnemyMovSpeed * sign(player.Pos.X-enemy.Pos.X)
	enemy.Vel.Y -= physics.Gravity * dt
	//enemy.Move(world, dt)

	playerIsCloseEnoughToStrike :=
		player.Pos.Sub(enemy.Pos).LengthSqr() <
			playerStrikeRange*playerStrikeRange

	stillAlive := !playerIsAttacking || !playerIsCloseEnoughToStrike
	return stillAlive
}

func (enemy *BasicEnemy) Draw(cam *camera.Camera) {
	camX := enemy.Pos.X - cam.X
	camY := enemy.Pos.Y - cam.Y

	spriteloader.DrawSpriteQuad(
		camX, camY,
		enemy.Size.X, enemy.Size.Y,
		basicEnemySpriteId,
	)
}
