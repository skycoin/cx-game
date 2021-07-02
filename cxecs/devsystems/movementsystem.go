package systems

import (
	"github.com/EngoEngine/ecs"
	components "github.com/skycoin/cx-game/cxecs/devcomponents"
	"github.com/skycoin/cx-game/cxecs/ecsconstants"
	"github.com/skycoin/cx-game/physics"
)

var (
	movAccumulator float32 = 0
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

func (ms *MovementSystem) Remove(entity ecs.BasicEntity) {
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

func (ms *MovementSystem) Priority() int {
	return ecsconstants.MOVEMENT_SYSTEM_PRIORITY
}

func (ms *MovementSystem) Update(dt float32) {
	// fmt.Println(len(ms.entities))
	movAccumulator += dt
	for movAccumulator >= physics.TimeStep {
		movAccumulator -= physics.TimeStep
		ms.fixedUpdate(physics.TimeStep)
	}
}

var mcounter int

func (ms *MovementSystem) fixedUpdate(dt float32) {
	for _, entity := range ms.entities {
		entity.Position = entity.Position.Add(entity.Velocity.Mult(dt))
	}
	mcounter += 1
	// fmt.Println(float64(mcounter)/glfw.GetTime(), "   movement")
	// fmt.Println("movement update")
}
