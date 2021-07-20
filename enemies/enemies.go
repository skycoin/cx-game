package enemies

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/spriteloader"
)

// TODO either create enemy as a side effect or return instance
func NewBasicEnemy(x, y float32) *agents.Agent {
	agent := agents.Agent {
		AgentType: constants.AGENT_ENEMY_MOB,
		AiHandlerID: constants.AI_HANDLER_WALK,
		DrawHandlerID: constants.DRAW_HANDLER_QUAD,
		PhysicsState: physics.Body{
			Size: cxmath.Vec2{X: 3.0, Y: 3.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		HealthComponent: agents.NewHealthComponent(5),
	}
	physics.RegisterBody(&agent.PhysicsState)
	return &agent
}

func NewLeapingEnemy(x,y float32) *agents.Agent {
	animation := spriteloader.NewSpriteAnimated("./assets/slime.json")
	action := animation.Action("Idle")
	agent := agents.Agent {
		AgentType: constants.AGENT_ENEMY_MOB,
		AiHandlerID: constants.AI_HANDLER_LEAP,
		DrawHandlerID: constants.DRAW_HANDLER_ANIM,
		PhysicsState: physics.Body{
			Size: cxmath.Vec2{X:2.0, Y: 2.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		HealthComponent: agents.NewHealthComponent(5),
		AnimationState: agents.AnimationState { Action: action },
	}
	physics.RegisterBody(&agent.PhysicsState)
	return &agent
}
