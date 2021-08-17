package agents

import (
	"log"

	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/render"
)

type AgentCreationOptions struct {
	X, Y float32
}

type AgentType struct {
	Name        string
	Category    constants.AgentCategory
	CreateAgent AgentCreator
}

func (at AgentType) Valid() bool {
	return at.Name != "" && at.CreateAgent != nil &&
		at.Category != constants.AGENT_CATEGORY_UNDEFINED
}

// TODO consider returning struct rather than pointer
type AgentCreator func(AgentCreationOptions) *Agent

var agentTypes [constants.NUM_AGENT_TYPES]AgentType

func init() {
	defer assertAllAgentTypesRegistered()
	RegisterAgentType(constants.AGENT_TYPE_SLIME, AgentType{
		Name:        "Slime",
		Category:    constants.AGENT_CATEGORY_ENEMY_MOB,
		CreateAgent: createSlime,
	})
	RegisterAgentType(constants.AGENT_TYPE_SPIDER_DRILL, AgentType{
		Name:        "Spider Drill",
		Category:    constants.AGENT_CATEGORY_ENEMY_MOB,
		CreateAgent: createSpiderDrill,
	})
	RegisterAgentType(constants.AGENT_TYPE_PLAYER, AgentType{
		Name:        "Player",
		Category:    constants.AGENT_CATEGORY_PLAYER,
		CreateAgent: createPlayer,
	})
}

func RegisterAgentType(id constants.AgentTypeID, agentType AgentType) {
	if !agentType.Valid() {
		log.Fatalf("invalid agent type %+v", agentType)
	}
	agentTypes[id] = agentType
}

func GetAgentType(id constants.AgentTypeID) AgentType {
	return agentTypes[id]
}

func createSlime(opts AgentCreationOptions) *Agent {
	x := opts.X
	y := opts.Y
	animation := anim.LoadAnimationFromJSON("./assets/slime.json")
	playback := animation.NewPlayback("Idle")
	agent := Agent{
		AgentCategory: constants.AGENT_CATEGORY_ENEMY_MOB,
		AiHandlerID:   constants.AI_HANDLER_LEAP,
		DrawHandlerID: constants.DRAW_HANDLER_ANIM,
		PhysicsState: physics.Body{
			Size: cxmath.Vec2{X: 2.0, Y: 2.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		HealthComponent:   NewHealthComponent(5),
		AnimationPlayback: playback,
	}
	//physics.RegisterBody(&agent.PhysicsState)
	return &agent
}

func createSpiderDrill(opts AgentCreationOptions) *Agent {
	x := opts.X
	y := opts.Y
	// TODO only load these once
	animation := anim.LoadAnimationFromJSON("./assets/spiderDrill.json")
	playback := animation.NewPlayback("Walk")
	agent := Agent{
		AgentCategory: constants.AGENT_CATEGORY_ENEMY_MOB,
		AiHandlerID:   constants.AI_HANDLER_DRILL,
		DrawHandlerID: constants.DRAW_HANDLER_ANIM,
		PhysicsState: physics.Body{
			Size: cxmath.Vec2{X: 2.0, Y: 2.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		HealthComponent:   NewHealthComponent(5),
		AnimationPlayback: playback,
	}
	//physics.RegisterBody(&agent.PhysicsState)
	return &agent
}

func createPlayer(opts AgentCreationOptions) *Agent {
	agent := Agent{
		AgentCategory: constants.AGENT_CATEGORY_PLAYER,
		AiHandlerID:   constants.AI_HANDLER_PLAYER,
		DrawHandlerID: constants.DRAW_HANDLER_PLAYER,
		PhysicsState: physics.Body{
			Pos:       cxmath.Vec2{X: opts.X, Y: opts.Y},
			Size:      cxmath.Vec2{X: 2.0 * constants.PLAYER_RENDER_TO_HITBOX, Y: 3},
			Direction: 1,
		},
		HealthComponent: NewHealthComponent(100),
		PlayerData: PlayerData{
			HelmetSpriteID: render.GetSpriteIDByName("helmet/1"),
			SuitSpriteID:   render.GetSpriteIDByName("suit:0"),
		},
	}
	//physics.RegisterBody(&agent.PhysicsState)
	return &agent
}

func assertAllAgentTypesRegistered() {
	for id, agentType := range agentTypes {
		if !agentType.Valid() {
			log.Fatalf("did not register agent type for id=%d", id)
		}
	}
}
