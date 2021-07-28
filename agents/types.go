package agents

import (
	"log"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/cxmath"
)

type AgentCreationOptions struct {
	X,Y float32
}

type AgentType struct {
	Name string
	Category constants.AgentCategory
	CreateAgent AgentCreator
}

func (at AgentType) Valid() bool {
	return at.Name!="" && at.CreateAgent!=nil &&
		at.Category!=constants.AGENT_CATEGORY_UNDEFINED
}

// TODO consider returning struct rather than pointer
type AgentCreator func(AgentCreationOptions) *Agent

var agentTypes [constants.NUM_AGENT_TYPES]AgentType

func init() {
	defer assertAllAgentTypesRegistered()
	RegisterAgentType(constants.AGENT_TYPE_SLIME, AgentType {
		Name: "Slime",
		Category: constants.AGENT_CATEGORY_ENEMY_MOB,
		CreateAgent: createSlime,
	})
	RegisterAgentType(constants.AGENT_TYPE_SPIDER_DRILL, AgentType {
		Name: "Spider Drill",
		Category: constants.AGENT_CATEGORY_ENEMY_MOB,
		CreateAgent: createSpiderDrill,
	})
}

func RegisterAgentType(id constants.AgentTypeID, agentType AgentType) {
	if !agentType.Valid() {
		log.Fatalf("invalid agent type %+v",agentType)
	}
	agentTypes[id] = agentType
}

func GetAgentType(id constants.AgentTypeID) AgentType {
	return agentTypes[id]
}

func createSlime(opts AgentCreationOptions) *Agent {
	x := opts.X; y := opts.Y;
	animation := spriteloader.NewSpriteAnimated("./assets/slime.json")
	action := animation.Action("Jump")
	agent := Agent{
		AgentCategory:     constants.AGENT_CATEGORY_ENEMY_MOB,
		AiHandlerID:   constants.AI_HANDLER_LEAP,
		DrawHandlerID: constants.DRAW_HANDLER_ANIM,
		PhysicsState: physics.Body{
			Size: cxmath.Vec2{X: 2.0, Y: 2.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		HealthComponent: NewHealthComponent(5),
		AnimationState:  AnimationState{Action: action},
	}
	physics.RegisterBody(&agent.PhysicsState)
	return &agent
}

func createSpiderDrill(opts AgentCreationOptions) *Agent {
	x := opts.X; y := opts.Y;
	spiderDrill := spriteloader.NewSpriteAnimated("./assets/spiderDrill.json")
	action := spiderDrill.Action("Walk")
	agent := Agent{
		AgentCategory:     constants.AGENT_CATEGORY_ENEMY_MOB,
		AiHandlerID:   constants.AI_HANDLER_DRILL,
		DrawHandlerID: constants.DRAW_HANDLER_ANIM,
		PhysicsState: physics.Body{
			Size: cxmath.Vec2{X: 2.0, Y: 2.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		HealthComponent: NewHealthComponent(5),
		AnimationState:  AnimationState{Action: action},
	}
	physics.RegisterBody(&agent.PhysicsState)
	return &agent
}

func assertAllAgentTypesRegistered() {
	for id,agentType := range agentTypes {
		if !agentType.Valid() {
			log.Fatalf("did not register agent type for id=%d",id)
		}
	}
}
