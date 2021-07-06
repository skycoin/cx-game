package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/skycoin/cx-game/cxecs/components"
)

type PlanetTile struct {
	ecs.BasicEntity
	components.TransformComponent
	components.RenderComponent
}
