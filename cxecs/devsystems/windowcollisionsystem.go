package systems

import (
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
	for _, entity := range wcs.entities {
		wcs.resolveCollision(entity)
	}
}

func (wcs *WindowCollisionSystem) resolveCollision(entity CollisionEntity) {

	if entity.Position.X-entity.Size.X <= -ecsconstants.GRID_WIDTH/2 { //left
		entity.Position.X = -10 + entity.Size.X
		entity.Velocity.X = -entity.Velocity.X

	} else if entity.Position.X+entity.Size.X >= ecsconstants.GRID_WIDTH/2 { //right
		entity.Position.X = 10 - entity.Size.X
		entity.Velocity.X = -entity.Velocity.X
	}

	if entity.Position.Y-entity.Size.Y <= -ecsconstants.GRID_HEIGHT/2 { //bottom
		entity.Position.Y = -7.5 + entity.Size.Y
		entity.Velocity.Y = -entity.Velocity.Y
	} else if entity.Position.Y+entity.Size.Y >= ecsconstants.GRID_HEIGHT/2 {
		entity.Position.Y = 7.5 - entity.Size.Y
		entity.Velocity.Y = -entity.Velocity.Y
	}

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
