package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/skycoin/cx-game/cxecs/components"
)

type MovementEntity struct {
	*ecs.BasicEntity
	*components.VelocityComponent
	*components.TransformComponent
}

type MovementSystem struct {
	entities []MovementEntity
}

func (*MovementSystem) New(*ecs.World) {

}

func (ms *MovementSystem) Add(entity *ecs.BasicEntity, velocityComponent *components.VelocityComponent, transformComponent *components.TransformComponent) {
	ms.entities = append(ms.entities, MovementEntity{entity, velocityComponent, transformComponent})
}

func (ms MovementSystem) Remove(entity ecs.BasicEntity) {
	delete := -1

	for index, e := range ms.entities {
		if e.ID() == entity.ID() {
			delete = index
			break
		}
	}
	if delete != -1 {
		ms.entities = append(ms.entities[:delete], ms.entities[delete+1:]...)
	}
}
func (ms MovementSystem) Update(dt float32) {
	// fmt.Println(len(ms.entities))
	for _, entity := range ms.entities {
		entity.Position = entity.Position.Add(entity.Velocity.Mult(dt))
	}
}
