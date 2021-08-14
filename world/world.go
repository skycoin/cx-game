package world

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/particles"
)

type Entities struct {
	Agents    agents.AgentList
	Particles particles.ParticleList
}

type World struct {
	Entities Entities
	Planet   Planet
}

func (world World) TileIsClear(layerID LayerID, x, y int) bool {
	layerTiles := world.Planet.GetLayerTiles(layerID)
	return world.Entities.Agents.TileIsClear(x, y) &&
		!world.Planet.TileExists(layerTiles, x, y)
}
