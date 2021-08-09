package agents

import (
	"log"

	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/render"
)

type Agent struct {
	AgentId           int
	AgentCategory     constants.AgentCategory
	AiHandlerID       types.AgentAiHandlerID
	PhysicsState      physics.Body
	PhysicsParameters physics.PhysicsParameters
	DrawHandlerID     types.AgentDrawHandlerID
	HealthComponent   HealthComponent
	TimeSinceDeath    float32
	WaitingFor        float32
	AnimationPlayback anim.Playback
	Direction         int
	InventoryID       types.InventoryID
	// only relevant to player agents - should probably refactor
	PlayerData
}

type PlayerData struct {
	SuitSpriteID   render.SpriteID
	HelmetSpriteID render.SpriteID
}

type HealthComponent struct {
	Current int
	Max     int
	Died    bool
}

func NewHealthComponent(max int) HealthComponent {
	return HealthComponent{Current: max, Max: max, Died: false}
}

func newAgent(id int) *Agent {
	agent := Agent{
		AgentCategory:     constants.AGENT_CATEGORY_UNDEFINED,
		AiHandlerID:       constants.AI_HANDLER_NULL,
		DrawHandlerID:     constants.DRAW_HANDLER_NULL,
		PhysicsState:      physics.Body{},
		PhysicsParameters: physics.PhysicsParameters{Radius: 5},
	}

	return &agent
}

func (a *Agent) FillDefaults() {
	// if have no direction, default to right-facing / no flip
	if a.PhysicsState.Direction == 0 {
		a.PhysicsState.Direction = 1
	}
}

//prefabs

func (a *Agent) SetPosition(x, y float32) {
	a.PhysicsState.Pos.X = x
	a.PhysicsState.Pos.Y = y
}

func (a *Agent) SetSize(x, y float32) {
	a.PhysicsState.Size.X = x
	a.PhysicsState.Size.Y = y
}

func (a *Agent) SetVelocity(x, y float32) {
	a.PhysicsState.Vel.X = x
	a.PhysicsState.Vel.Y = y
}

func (a *Agent) TakeDamage(amount int) {
	a.HealthComponent.Current -= amount
	if a.HealthComponent.Current <= 0 {
		a.HealthComponent.Died = true
	}
}

func (a *Agent) Died() bool {
	return a.HealthComponent.Died
}

func (a *Agent) IsWaiting() bool {
	return a.WaitingFor > 0
}

func (a *Agent) WaitFor(seconds float32) {
	a.WaitingFor = seconds
}

func (a *Agent) Validate() {
	if a.AgentCategory == constants.AGENT_CATEGORY_UNDEFINED {
		log.Fatalf("Cannot create agent with undefined category: %+v", a)
	}
}
