package cxecs

import (
	"github.com/EngoEngine/ecs"
	"github.com/skycoin/cx-game/cxecs/entities"
	"github.com/skycoin/cx-game/cxecs/systems"
	"github.com/skycoin/cx-game/spriteloader"
)

var world ecs.World

func InitDev() {
	world = ecs.World{}
	loadAssets()
	world.AddSystem(&systems.RenderSystem{})
	world.AddSystem(&systems.MovementSystem{})
	world.AddSystem(&systems.CollisionSystem{})
	// world.AddSystem(systems.SimpleSystem{})

	var enemies []*entities.Enemy

	for i := 0; i < 5; i++ {
		enemies = append(enemies, entities.NewEnemy())
	}

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *systems.RenderSystem:

		case *systems.MovementSystem:
			for _, enemy := range enemies {
				sys.Add(&enemy.BasicEntity, &enemy.VelocityComponent, &enemy.TransformComponent)
			}
		case *systems.CollisionSystem:
			{
				for _, enemy := range enemies {
					sys.Add(&enemy.BasicEntity, &enemy.VelocityComponent, &enemy.TransformComponent)
				}
			}
		}

	}
}

func loadAssets() {
	spriteloader.LoadSingleSprite("./assets/enemies/basic-enemy.png", "enemy")
}

func Update(dt float32) {
	world.Update(dt)
}
