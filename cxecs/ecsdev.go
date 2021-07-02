package cxecs

import (
	"github.com/EngoEngine/ecs"
	entities "github.com/skycoin/cx-game/cxecs/deventities"
	systems "github.com/skycoin/cx-game/cxecs/devsystems"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/spriteloader"
)

var (
	worldDev ecs.World
	window   *render.Window
)

func InitDev(window *render.Window) {
	worldDev = ecs.World{}
	loadAssetsDev()
	worldDev.AddSystem(&systems.RenderSystem{})
	worldDev.AddSystem(&systems.MovementSystem{})
	worldDev.AddSystem(&systems.CollisionSystem{})
	worldDev.AddSystem(&systems.WindowCollisionSystem{})
	// world.AddSystem(systems.SimpleSystem{})

	var enemies []*entities.Enemy

	for i := 0; i < 5; i++ {
		enemies = append(enemies, entities.NewEnemy())
	}

	for _, system := range worldDev.Systems() {
		switch sys := system.(type) {
		case *systems.RenderSystem:
			for _, enemy := range enemies {
				sys.Add(&enemy.BasicEntity, &enemy.RenderComponent, &enemy.TransformComponent)
			}
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
		case *systems.WindowCollisionSystem:
			for _, enemy := range enemies {
				sys.Add(&enemy.BasicEntity, &enemy.VelocityComponent, &enemy.TransformComponent)
			}
		}

	}
}

func loadAssetsDev() {
	spriteloader.LoadSingleSprite("./assets/enemies/basic-enemy.png", "enemy")
}

func UpdateDev(dt float32) {
	worldDev.Update(dt)
}

func GetWindowHandle() *render.Window {
	return window
}
