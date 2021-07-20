package constants

type AgentType int

const (
	MAX_AGENTS int = 64

	AGENT_UNDEFINED AgentType = iota
	AGENT_PLAYER
	AGENT_FRIENDLY_MOB
	AGENT_ENEMY_MOB
)
