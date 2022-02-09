package agents

import (
	"log"

	"github.com/skycoin/cx-game/components/types"
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
	Category    types.AgentCategory
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
	RegisterAgentType(constants.AGENT_TYPE_ENEMY_FLOATING, AgentType {
		Name:        "Floater",
		Category:    constants.AGENT_CATEGORY_ENEMY_MOB,
		CreateAgent: createFloatingEnemy,
	})
	RegisterAgentType(constants.AGENT_TYPE_SPIDER_DRILL, AgentType{
		Name:        "Spider Drill",
		Category:    constants.AGENT_CATEGORY_ENEMY_MOB,
		CreateAgent: createSpiderDrill,
	})
	RegisterAgentType(constants.AGENT_TYPE_GRASS_HOPPER, AgentType{
		Name:        "Grass Hopper",
		Category:    constants.AGENT_CATEGORY_ENEMY_MOB,
		CreateAgent: createGrassHopper,
	})
	RegisterAgentType(constants.AGENT_TYPE_ENEMY_SOLDIER, AgentType{
		Name:        "Enemy Soldier",
		Category:    constants.AGENT_CATEGORY_ENEMY_MOB,
		CreateAgent: createEnemySoldier,
	})
	RegisterAgentType(constants.AGENT_TYPE_PLAYER, AgentType{
		Name:        "Player",
		Category:    constants.AGENT_CATEGORY_PLAYER,
		CreateAgent: createPlayer,
	})
	// RegisterAgentType(constants.AGENT_TYPE_PLAYER, AgentType{
	// 	Name:        "Player",
	// 	Category:    constants.AGENT_CATEGORY_PLAYER,
	// 	CreateAgent: createSkeletonPlayer,
	// })
}

func RegisterAgentType(id types.AgentTypeID, agentType AgentType) {
	if !agentType.Valid() {
		log.Fatalf("invalid agent type %+v", agentType)
	}
	agentTypes[id] = agentType
}

func GetAgentType(id types.AgentTypeID) AgentType {
	return agentTypes[id]
}

func createSlime(opts AgentCreationOptions) *Agent {
	x := opts.X
	y := opts.Y
	animation := anim.LoadAnimationFromJSON("./assets/slime.json")
	playback := animation.NewPlayback("Idle")
	agent := Agent{
		Handlers: AgentHandlers{
			AI:   constants.AI_HANDLER_LEAP,
			Draw: constants.DRAW_HANDLER_ANIM,
		},
		Meta: AgentMeta{
			Category: constants.AGENT_CATEGORY_ENEMY_MOB,
			Type:     constants.AGENT_TYPE_SLIME,
			PhysicsParameters: physics.PhysicsParameters{
				Radius: 1,
			},
		},
		Transform: physics.Body{
			Size: cxmath.Vec2{X: 2.0, Y: 2.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		Health:            NewHealthComponent(constants.HEALTH_SLIME),
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
		Handlers: AgentHandlers{
			AI:   constants.AI_HANDLER_DRILL,
			Draw: constants.DRAW_HANDLER_ANIM,
		},
		Meta: AgentMeta{
			Category: constants.AGENT_CATEGORY_ENEMY_MOB,
			Type:     constants.AGENT_TYPE_SPIDER_DRILL,
		},
		Transform: physics.Body{
			Size: cxmath.Vec2{X: 2.0, Y: 2.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		Health:            NewHealthComponent(constants.HEALTH_SPIDERDRILL),
		AnimationPlayback: playback,
	}
	//physics.RegisterBody(&agent.PhysicsState)
	return &agent
}

func createGrassHopper(opts AgentCreationOptions) *Agent {
	x := opts.X
	y := opts.Y
	animation := anim.LoadAnimationFromJSON("./assets/grassHopper.json")
	playback := animation.NewPlayback("Idle")
	agent := Agent{
		Handlers: AgentHandlers{
			AI:   constants.AI_HANDLER_GRASSHOPPER,
			Draw: constants.DRAW_HANDLER_ANIM,
		},
		Meta: AgentMeta{
			Category: constants.AGENT_CATEGORY_ENEMY_MOB,
			Type:     constants.AGENT_TYPE_GRASS_HOPPER,
		},
		Transform: physics.Body{
			Size: cxmath.Vec2{X: 2.0, Y: 3.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		Health:            NewHealthComponent(constants.HEALTH_GRASSHOPPER),
		AnimationPlayback: playback,
	}

	return &agent
}

func createEnemySoldier(opts AgentCreationOptions) *Agent {
	x := opts.X
	y := opts.Y
	animation := anim.LoadAnimationFromJSON("./assets/enemy_soldier.json")
	playback := animation.NewPlayback("Idle")
	agent := Agent{
		Handlers: AgentHandlers{
			AI:   constants.AI_HANDLER_ENEMYSOLDIER,
			Draw: constants.DRAW_HANDLER_ANIM,
		},
		Meta: AgentMeta{
			Category: constants.AGENT_CATEGORY_ENEMY_MOB,
			Type:     constants.AGENT_TYPE_ENEMY_SOLDIER,
		},
		Transform: physics.Body{
			Size: cxmath.Vec2{X: 2.0, Y: 3.0},
			Pos:  cxmath.Vec2{X: x, Y: y},
		},
		Health:            NewHealthComponent(constants.HEALTH_ENEMYSOLDIER),
		AnimationPlayback: playback,
	}

	return &agent
}

func createFloatingEnemy(opts AgentCreationOptions) *Agent {
	agent := Agent{
		Handlers: AgentHandlers{
			AI:   constants.AI_HANDLER_FLOATING,
			Draw: constants.DRAW_HANDLER_COLOR,
		},
		Meta: AgentMeta{
			Category: constants.AGENT_CATEGORY_ENEMY_MOB,
			Type:     constants.AGENT_TYPE_ENEMY_FLOATING,
			PhysicsParameters: physics.PhysicsParameters{
				Radius: 1,
			},
		},
		Transform: physics.Body{
			Size:           cxmath.Vec2{X: 2.0, Y: 2.0},
			Pos:            cxmath.Vec2{X: opts.X, Y: opts.Y},
			IgnoresGravity: true,
		},
		Health: NewHealthComponent(constants.HEALTH_ENEMYFLOATING),
	}
	return &agent
}

func createPlayer(opts AgentCreationOptions) *Agent {
	agent := Agent{
		Handlers: AgentHandlers{
			AI:   constants.AI_HANDLER_PLAYER,
			Draw: constants.DRAW_HANDLER_PLAYER,
		},

		Meta: AgentMeta{
			Category: constants.AGENT_CATEGORY_PLAYER,
			Type:     constants.AGENT_TYPE_PLAYER,
			PlayerData: PlayerData{
				HelmetSpriteID: render.GetSpriteIDByName("helmet/1"),
				SuitSpriteID:   render.GetSpriteIDByName("suit:0"),
			},
		},
		Transform: physics.Body{
			Pos:       cxmath.Vec2{X: opts.X, Y: opts.Y},
			Size:      cxmath.Vec2{X: 2.0 * constants.PLAYER_RENDER_TO_HITBOX, Y: 3},
			Direction: 1,
		},
		Health: NewHealthComponent(constants.HEALTH_PLAYER),
	}
	//physics.RegisterBody(&agent.PhysicsState)
	return &agent
}

func createSkeletonPlayer(opts AgentCreationOptions) *Agent {
	agent := Agent{
		Handlers: AgentHandlers{
			AI:   constants.AI_HANDLER_PLAYER,
			Draw: constants.DRAW_HANDLER_PLAYER,
		},
		Meta: AgentMeta{
			Category:   constants.AGENT_CATEGORY_PLAYER,
			Type:       constants.AGENT_TYPE_PLAYER,
			PlayerData: PlayerData{
				// HelmetSpriteID: render.GetSpriteIDByName("helmet/1"),
				// SuitSpriteID:   render.GetSpriteIDByName("suit:0"),
			},
		},
		Transform: physics.Body{
			Pos:       cxmath.Vec2{X: opts.X, Y: opts.Y},
			Size:      cxmath.Vec2{X: 2.0 * constants.PLAYER_RENDER_TO_HITBOX, Y: 3},
			Direction: 1,
		},
		Health: NewHealthComponent(constants.HEALTH_PLAYER),
	}
	return &agent
}

func assertAllAgentTypesRegistered() {
	for id, agentType := range agentTypes {
		if !agentType.Valid() {
			log.Fatalf("did not register agent type for id=%d", id)
		}
	}
}
