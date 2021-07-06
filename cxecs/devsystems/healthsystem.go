package systems

import (
	"github.com/EngoEngine/ecs"
	components "github.com/skycoin/cx-game/cxecs/devcomponents"
	"github.com/skycoin/cx-game/cxecs/ecsconstants"
	"github.com/skycoin/cx-game/physics"
)

var (
	healthAccumulator float32 = 0
)

type HealthEntity struct {
	*ecs.BasicEntity
	*components.HealthComponent
	*components.TransformComponent
}

type HealthSystem struct {
	entities []HealthEntity
}

func (*HealthSystem) New(*ecs.World) {

}

func (hs *HealthSystem) Add(entity *ecs.BasicEntity, healthComponent *components.HealthComponent, transformComponent *components.TransformComponent) {
	hs.entities = append(hs.entities, HealthEntity{entity, healthComponent, transformComponent})
}

func (hs *HealthSystem) Remove(entity ecs.BasicEntity) {
	delete := -1

	for index, e := range hs.entities {
		if e.ID() == entity.ID() {
			delete = index
			break
		}
	}
	if delete != -1 {
		hs.entities = append(hs.entities[:delete], hs.entities[delete+1:]...)
	}
}

func (hs *HealthSystem) Priority() int {
	return ecsconstants.HEALTH_SYSTEM_PRIORITY
}

func (hs *HealthSystem) Update(dt float32) {
	// fmt.Println(len(ms.entities))
	movAccumulator += dt
	for movAccumulator >= physics.TimeStep {
		movAccumulator -= physics.TimeStep
		hs.fixedUpdate(physics.TimeStep)
	}
}

func (hs *HealthSystem) fixedUpdate(dt float32) {
	//todo
	// for _, entity := range hs.entities {
	// 	entity
	// }
}
