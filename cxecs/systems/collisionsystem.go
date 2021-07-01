package systems

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/skycoin/cx-game/cxecs/components"
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

func (rs *CollisionSystem) Add(entity *ecs.BasicEntity, velocityComponent *components.VelocityComponent, transformComponent *components.TransformComponent) {
	rs.entities = append(rs.entities, CollisionEntity{entity, velocityComponent, transformComponent})
}

func (rs CollisionSystem) Remove(entity ecs.BasicEntity) {
	delete := -1

	for index, e := range rs.entities {
		if e.ID() == entity.ID() {
			delete = index
			break
		}
	}
	if delete != -1 {
		rs.entities = append(rs.entities[:delete], rs.entities[delete+1:]...)
	}
}
func (rs CollisionSystem) Update(dt float32) {
	// fmt.Println(len(rs.entities))
	for _, entity1 := range rs.entities {
		for _, entity2 := range rs.entities {
			if entity1.ID() != entity2.ID() {
				if collides(entity1, entity2) {
					resolve(entity1, entity2)
				}
			}
		}
	}
}

func collides(entity1, entity2 CollisionEntity) bool {
	//aabb
	// if entity1.Position.X + entity1.Size.X

	if entity1.Position.X+entity1.Size.X/2 < entity2.Position.X-entity2.Size.X/2 &&
		entity1.Position.Y-entity1.Size.Y/2 < entity2.Position.Y-entity1.Size.Y/2 {
		return true
	}
	return false
}

func resolve(entity1, entity2 CollisionEntity) {
	fmt.Println("collided!")
}
