package entities

import (
	"github.com/EngoEngine/ecs"
	"github.com/skycoin/cx-game/cxecs/components"
)

type Player struct {
	ecs.BasicEntity
	components.TransformComponent
	components.RenderComponent
}
