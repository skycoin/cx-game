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
	Tick int
	Entities Entities
	Planet   Planet
	Stats    WorldStats
}

func (world World) TileIsClear(layerID LayerID, x, y int) bool {
	layerTiles := world.Planet.GetLayerTiles(layerID)

	tileClear := !world.Planet.TileExists(layerTiles, x, y)
	agentsClear := world.Entities.Agents.TileIsClear(x, y)
	isClear :=
		 tileClear &&
		( layerID != TopLayer || agentsClear )
	return isClear
}
