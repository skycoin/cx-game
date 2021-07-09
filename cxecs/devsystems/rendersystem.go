package systems

import (
	"github.com/EngoEngine/ecs"
	components "github.com/skycoin/cx-game/cxecs/devcomponents"
	"github.com/skycoin/cx-game/cxecs/ecsconstants"
	"github.com/skycoin/cx-game/spriteloader"
)

type RenderEntity struct {
	*ecs.BasicEntity
	*components.RenderComponent
	*components.TransformComponent
}

type RenderSystem struct {
	entities []RenderEntity
}

func (*RenderSystem) New(*ecs.World) {

}

func (rs *RenderSystem) Add(entity *ecs.BasicEntity, renderComponent *components.RenderComponent, transformComponent *components.TransformComponent) {
	rs.entities = append(rs.entities, RenderEntity{entity, renderComponent, transformComponent})
}

func (rs *RenderSystem) Remove(entity ecs.BasicEntity) {
	delete := -1

	for index, e := range rs.entities {
		if e.ID() == entity.ID() {
			delete = index
			break
		}
	}
	if delete != -1 {
		rs.entities = append(rs.entities[:delete], rs.entities[delete+1:]...)
	}
}

func (rs *RenderSystem) Priority() int {
	return ecsconstants.RENDER_SYSTEM_PRIORITY
}

func (rs *RenderSystem) Update(dt float32) {
	// fmt.Println(len(rs.entities))
	for _, entity := range rs.entities {
		spriteloader.DrawSpriteQuad(
			entity.Position.X,
			entity.Position.Y,
			entity.Size.X,
			entity.Size.Y,
			(entity.SpriteId),
		)
	}
	// fmt.Println("render system")
}
