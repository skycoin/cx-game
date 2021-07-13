package constants

type AgentType int

const (
	MAX_AGENTS = 64

	AGENT_PLAYER = iota
	AGENT_NEUTRAL_MOB
	AGENT_FRIENDLY_MOB
	AGENT_ENEMY_MOB
)
