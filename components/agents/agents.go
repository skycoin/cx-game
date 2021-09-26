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
	Meta              AgentMeta
	Handlers          AgentHandlers
	Timers            AgentTimers
	AgentId           types.AgentID
	Inventory         types.InventoryID
	PhysicsState      physics.Body
	HealthComponent   HealthComponent
	AnimationPlayback anim.Playback
	Direction         int
	InventoryID       types.InventoryID
	// only relevant to player agents - should probably refactor
	PlayerData
}

type AgentMeta struct {
	Category          constants.AgentCategory
	Type              constants.AgentTypeID
	PhysicsParameters physics.PhysicsParameters
}

type AgentHandlers struct {
	Draw types.AgentDrawHandlerID
	AI   types.AgentAiHandlerID
}

type AgentTimers struct {
	TimeSinceDeath float32
	WaitingFor     float32
}
type PlayerData struct {
	SuitSpriteID         render.SpriteID
	HelmetSpriteID       render.SpriteID
	IgnoringPlatformsFor float32
}

type HealthComponent struct {
	Current int
	Max     int
	Died    bool
}

func NewHealthComponent(max int) HealthComponent {
	return HealthComponent{Current: max, Max: max, Died: false}
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
