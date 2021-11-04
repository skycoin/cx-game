package worldimport

import (
	"log"
	"time"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/world"
	"github.com/lafriks/go-tiled"
)

func ImportWorld(tmxPath string) world.World {
	start := time.Now()
	tiledMap, err := tiled.LoadFromFile(tmxPath)
	if err != nil {
		log.Fatalf("import world: %v", err)
	}
	elapsedTiledLoad := time.Since(start)
	log.Printf("load %s took %s", tmxPath, elapsedTiledLoad)
	planet := world.NewPlanet(int32(tiledMap.Width), int32(tiledMap.Height))
	for _, tiledLayer := range tiledMap.Layers {
		layerID, foundLayerID := world.LayerIDForName(tiledLayer.Name)
		if foundLayerID {
			importLayer(planet, tiledLayer, tmxPath, layerID)
		}
	}
	return world.World{
		Planet: *planet,
		Entities: world.Entities{
			Agents: *agents.NewAgentList(),
		},
		Stats:     world.NewWorldStats(),
		TimeState: world.NewTimeState(),
	}
}
