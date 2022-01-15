package constants

import "github.com/skycoin/cx-game/components/types"

const (
	MAX_AGENTS int = 64
	// how high agents are teleported for vertical wrap around
	// and also how high agents can go
	HEIGHT_LIMIT float32 = 256
)

const (
	AGENT_CATEGORY_UNDEFINED types.AgentCategory = iota
	AGENT_CATEGORY_PLAYER
	AGENT_CATEGORY_FRIENDLY_MOB
	AGENT_CATEGORY_ENEMY_MOB

	AGENT_CATEGORY_SPINE
)

const (
	AGENT_TYPE_SLIME types.AgentTypeID = iota
	AGENT_TYPE_SPIDER_DRILL
	AGENT_TYPE_GRASS_HOPPER
	AGENT_TYPE_ENEMY_SOLDIER
	AGENT_TYPE_PLAYER

	//test spine
	AGENT_TYPE_SPINE_TEST

	NUM_AGENT_TYPES // DO NOT SET THIS MANUALLY
)

const (
	PLAYER_RENDER_TO_HITBOX float32 = 43.0 / 64
)
