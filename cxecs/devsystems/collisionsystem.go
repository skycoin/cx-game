package systems

import (
	"math/rand"

	"github.com/EngoEngine/ecs"
	components "github.com/skycoin/cx-game/cxecs/devcomponents"
	"github.com/skycoin/cx-game/cxecs/ecsconstants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/physics"
)

var (
	colAccumulator float32 = 0
)

type CollisionEntity struct {
	*ecs.BasicEntity
	*components.VelocityComponent
	*components.TransformComponent
}

type CollisionSystem struct {
	entities []CollisionEntity
}

func (*CollisionSystem) New(*ecs.World) {

}

func (cs *CollisionSystem) Add(entity *ecs.BasicEntity, velocityComponent *components.VelocityComponent, transformComponent *components.TransformComponent) {
	cs.entities = append(cs.entities, CollisionEntity{entity, velocityComponent, transformComponent})
}

func (cs *CollisionSystem) Remove(entity ecs.BasicEntity) {
	delete := -1

	for index, e := range cs.entities {
		if e.ID() == entity.ID() {
			delete = index
			break
		}
	}
	if delete != -1 {
		cs.entities = append(cs.entities[:delete], cs.entities[delete+1:]...)
	}
}

func (cs *CollisionSystem) Priority() int {
	return ecsconstants.COLLISION_SYSTEM_PRIORITY
}
func (cs *CollisionSystem) Update(dt float32) {
	//fixed tick
	colAccumulator += dt

	for colAccumulator >= physics.TimeStep {
		colAccumulator -= physics.TimeStep
		cs.fixedUpdate()
	}
	//todo quadtree

}

var ccounter int

func (cs *CollisionSystem) fixedUpdate() {
	// n^2
	for _, entity1 := range cs.entities {
		for _, entity2 := range cs.entities {
			if entity1.ID() != entity2.ID() {
				if cs.collides(entity1, entity2) {
					cs.resolve(entity1, entity2)
				}
			}
		}
	}
	ccounter += 1
	// fmt.Println(float64(ccounter)/glfw.GetTime(), "    csystem")
	// fmt.Println("physics")
}
func (cs *CollisionSystem) collides(entity1, entity2 CollisionEntity) bool {
	//aabb collision detection

	// if player1.left < player2.right &&
	// 	  player1.right > player2.left &&
	// 	  player1.bottom < player2.top &&
	//	  player1.top > player2.bottom

	// https://learnopengl.com/img/in-practice/breakout/collisions_overlap.png

	if entity1.Position.X-entity1.Size.X/2 < entity2.Position.X+entity2.Size.X/2 &&
		entity1.Position.X+entity1.Size.X/2 > entity2.Position.X &&
		entity1.Position.Y-entity1.Size.Y/2 < entity2.Position.Y+entity2.Size.Y/2 &&
		entity1.Position.Y+entity1.Size.Y/2 > entity2.Position.Y-entity1.Size.Y/2 {
		return true
	}
	return false
}

var eps float32 = 0.1

func (cs *CollisionSystem) resolve(entity1, entity2 CollisionEntity) {

	// entity1VelocityOld := entity1.Velocity
	// entity1.Velocity = entity1.Velocity.Mult(-1).Add(entity2.Velocity.Normalize())
	// entity2.Velocity = entity2.Velocity.Mult(-1).Add(entity1VelocityOld.Normalize())
	entity1.Velocity = entity1.Velocity.Mult(-1 * rand.Float32() * 2)
	entity2.Velocity = entity2.Velocity.Mult(-1 * rand.Float32() * 2)
	entity1.Position.X += cxmath.Sign(entity1.Velocity.X) * eps
	entity1.Position.Y += cxmath.Sign(entity1.Velocity.Y) * eps
	entity2.Position.X += cxmath.Sign(entity2.Velocity.X) * eps
	entity2.Position.Y += cxmath.Sign(entity2.Velocity.Y) * eps
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
