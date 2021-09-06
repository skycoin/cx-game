package worldgen

import (
	"github.com/skycoin/cx-game/world/mapgen"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/components/agents"
)

func GenerateWorld() world.World {
	return world.World{
		Entities: world.Entities{
			Agents: *agents.NewAgentList(),
		},
		Planet: *mapgen.GeneratePlanet(),
		Stats: world.NewWorldStats(),
	}
}
