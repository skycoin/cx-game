package world

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/components/particles"
)

type Entities struct {
	Agents agents.AgentList
	Particles particles.ParticleList
}

type World struct {
	Entities Entities
	Planet Planet
}
