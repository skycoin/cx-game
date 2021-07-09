package pathfinding

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/physics"
)

type BehaviourContext struct {
	Self *physics.Body
	PlayerPos mgl32.Vec2
}

type Behaviour interface {
	Follow(BehaviourContext)
}
type BehaviourID uint32

var behaviours = []Behaviour{}

func (id BehaviourID) Get() Behaviour {
	return behaviours[id]
}

func RegisterBehaviour(behaviour Behaviour) BehaviourID {
	id := BehaviourID(len(behaviours))
	behaviours = append(behaviours,behaviour)
	return id
}
