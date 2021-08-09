package movement_state

import (
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/physics/movement"
)

type NormalState struct {
}

func (state *NormalState) Enter() {

}

func (state *NormalState) HandleInput(agent *models.Player) {
	if input.GetButtonDown("fly") && agent.TryChangeMovementState(movement.FLYING) {

	}
}

func (state *NormalState) Update(agent *models.Player) {
	if agent.Collisions.Below {
		agent.ResetJumpCounter()
	} else if agent.Collisions.Left || agent.Collisions.Right && agent.IsMovementTypePresent(movement.WALL_SLIDING) {
		agent.TryChangeMovementState(movement.WALL_SLIDING)
	}
}

func (state *NormalState) Exit() {

}
