package agents

import (
	"log"

	"github.com/skycoin/cx-game/common"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/input/inputhandler"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/physics"
)

type Agent struct {
	PlayerControlled  bool
	AgentId           types.AgentID
	Meta              AgentMeta
	Handlers          AgentHandlers
	Timers            AgentTimers
	Transform         physics.Body
	Health            HealthComponent
	AnimationPlayback anim.Playback
	InventoryID       types.InventoryID
	CS                common.Bitset
}

type AgentHandlers struct {
	Draw types.AgentDrawHandlerID
	AI   types.AgentAiHandlerID
}

type AgentTimers struct {
	TimeSinceDeath float32
	WaitingFor     float32
}

// func newAgent(id int) *Agent {
// 	// agent := Agent{
// 	// 	AgentCategory:     constants.AGENT_CATEGORY_UNDEFINED,
// 	// 	AiHandlerID:       constants.AI_HANDLER_NULL,
// 	// 	DrawHandlerID:     constants.DRAW_HANDLER_NULL,
// 	// 	PhysicsState:      physics.Body{},
// 	// 	PhysicsParameters: physics.PhysicsParameters{Radius: 5},
// 	// }

// 	return &agent
// }

func (a *Agent) FillDefaults() {
	// if have no direction, default to right-facing / no flip
	if a.Transform.Direction == 0 {
		a.Transform.Direction = 1
	}
}

//prefabs

func (a *Agent) SetPosition(x, y float32) {
	a.Transform.Pos.X = x
	a.Transform.Pos.Y = y
}

func (a *Agent) SetSize(x, y float32) {
	a.Transform.Size.X = x
	a.Transform.Size.Y = y
}

func (a *Agent) SetVelocity(x, y float32) {
	a.Transform.Vel.X = x
	a.Transform.Vel.Y = y
}

func (a *Agent) TakeDamage(amount int) {
	a.Health.Current -= amount
	if a.Health.Current <= 0 {
		a.Health.Died = true
	}
}

func (a *Agent) Died() bool {
	return a.Health.Died
}

func (a *Agent) IsWaiting() bool {
	return a.Timers.WaitingFor > 0
}

func (a *Agent) WaitFor(seconds float32) {
	a.Timers.WaitingFor = seconds
}

func (a *Agent) Validate() {
	if a.Meta.Category == constants.AGENT_CATEGORY_UNDEFINED {
		log.Fatalf("Cannot create agent with undefined category: %+v", a)
	}
}

func (a *Agent) SetControlState(cs common.Bitset) {

}

func (a *Agent) ApplyControlState() {
	moveLeft := a.CS.GetBit(inputhandler.MOVE_LEFT)
	moveRight := a.CS.GetBit(inputhandler.MOVE_RIGHT)
	jump := a.CS.GetBit(inputhandler.JUMP)
	crouch := a.CS.GetBit(inputhandler.CROUCH)

	
}
