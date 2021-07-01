package systems

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	components "github.com/skycoin/cx-game/cxecs/devcomponents"
	"github.com/skycoin/cx-game/cxecs/ecsconstants"
	"github.com/skycoin/cx-game/physics"
)

var (
	wcolAccumulator float32 = 0
)

type WindowCollisionEntity struct {
	*ecs.BasicEntity
	*components.VelocityComponent
	*components.TransformComponent
}

type WindowCollisionSystem struct {
	entities []CollisionEntity
}

func (*WindowCollisionSystem) New(*ecs.World) {

}

func (wcs *WindowCollisionSystem) Add(entity *ecs.BasicEntity, velocityComponent *components.VelocityComponent, transformComponent *components.TransformComponent) {
	wcs.entities = append(wcs.entities, CollisionEntity{entity, velocityComponent, transformComponent})
}

func (wcs *WindowCollisionSystem) Remove(entity ecs.BasicEntity) {
	delete := -1

	for index, e := range wcs.entities {
		if e.ID() == entity.ID() {
			delete = index
			break
		}
	}
	if delete != -1 {
		wcs.entities = append(wcs.entities[:delete], wcs.entities[delete+1:]...)
	}
}

func (wcs *WindowCollisionSystem) Priority() int {
	return ecsconstants.COLLISION_SYSTEM_PRIORITY
}
func (wcs *WindowCollisionSystem) Update(dt float32) {
	//fixed tick
	colAccumulator += dt

	for colAccumulator >= physics.TimeStep {
		colAccumulator -= physics.TimeStep
		wcs.fixedUpdate()
	}
	//todo quadtree

}

func (wcs *WindowCollisionSystem) fixedUpdate() {
	// n^2
	for _, entity1 := range wcs.entities {
		for _, entity2 := range wcs.entities {
			if entity1.ID() != entity2.ID() {
				if wcs.collides(entity1, entity2) {
					wcs.resolve(entity1, entity2)
				}
			}
		}
	}
}
func (wcs *WindowCollisionSystem) collides(entity1, entity2 CollisionEntity) bool {
	//aabb collision detection

	// if player1.left < player2.right &&
	// 	  player1.right > player2.left &&
	// 	  player1.bottom < player2.top &&
	//	  player1.top > player2.bottom

	// https://learnopengl.com/img/in-practice/breakout/collisions_overlap.png

	if entity1.Position.X-entity1.Size.X/2 < entity2.Position.X+entity2.Size.X/2 &&
		entity1.Position.X+entity1.Size.X/2 > entity2.Position.X &&
		entity1.Position.Y-entity1.Size.Y/2 < entity2.Position.Y+entity2.Size.Y &&
		entity1.Position.Y+entity1.Size.Y/2 > entity2.Position.Y-entity1.Size.Y/2 {
		return true
	}
	return false
}

func (wcs *WindowCollisionSystem) resolve(entity1, entity2 CollisionEntity) {
	fmt.Println("collided!")
}

/*
if(player1.x < player2.x + player2.width &&
    player1.x + player1.width > player2.x &&
    player1.y < player2.y + player2.height &&
    player1.y + player1.height > player2.y)
{
    System.out.println("Collision Detected");
}
*/
