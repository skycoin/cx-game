package worldgen

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/procgen/mapgen"
)

func GenerateWorld() world.World {
	return world.World{
		Entities: world.Entities{
			Agents: *agents.NewAgentList(),
		},
		Planet:    *mapgen.GeneratePlanet(),
		Stats:     world.NewWorldStats(),
		TimeState: world.NewTimeState(),
	}
}
