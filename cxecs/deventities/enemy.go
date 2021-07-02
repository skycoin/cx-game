package entities

import (
	"math/rand"

	"github.com/EngoEngine/ecs"
	components "github.com/skycoin/cx-game/cxecs/devcomponents"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
)

var (
	movSpeed float32 = 5
)

type Enemy struct {
	ecs.BasicEntity
	components.RenderComponent
	components.TransformComponent
	components.VelocityComponent
	components.HealthComponent
}

func NewEnemy() *Enemy {
	enemy := &Enemy{ecs.NewBasic(),
		components.RenderComponent{},
		components.TransformComponent{},
		components.VelocityComponent{},
		components.HealthComponent{},
	}

	enemy.RenderComponent.SpriteId = spriteloader.GetSpriteIdByName("enemy")

	enemy.TransformComponent.Position = cxmath.Vec2{float32(rand.Intn(10) - 5), float32(rand.Intn(10) - 5)}

	enemy.TransformComponent.Size = cxmath.Vec2{1, 1}

	xVel := (rand.Float32() - 0.5) * movSpeed
	yVel := (rand.Float32() - 0.5) * movSpeed
	enemy.Velocity = cxmath.Vec2{xVel, yVel}

	enemy.Health_Max = rand.Intn(10) + 20
	enemy.Health = enemy.Health_Max

	return enemy

}
