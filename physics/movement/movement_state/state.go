package movement_state

import (
	"github.com/skycoin/cx-game/input"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/physics/movement"
)

type State interface {
	Enter(*models.Player)
	HandleInput(*models.Player)
	Update(*models.Player)
	Exit(*models.Player)
}

//normal

//flying
type FlyingState struct {
}

func (state *FlyingState) Enter(agent *models.Player) {

}

func (state *FlyingState) HandleInput(agent *models.Player) {
	if input.GetButtonDown("fly") {
		agent.TryChangeMovementState(movement.NORMAL)
	}
}
