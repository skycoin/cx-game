// a world item is an item that is "floating" in the world
package item;

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/cxmath"
)

const worldItemSize = 0.5
type WorldItem struct {
	physics.Body
	ItemTypeId uint32
}

func NewWorldItem(ItemTypeId uint32) WorldItem {
	return WorldItem {
		Body: physics.Body {
			Size: physics.Vec2 { X: worldItemSize, Y: worldItemSize },
		},
		ItemTypeId: ItemTypeId,
	}
}

func (item WorldItem) Draw(cam *camera.Camera) {
	spriteId := itemTypes[item.ItemTypeId].SpriteID
	z := -spriteloader.SpriteRenderDistance
	view := mgl32.Translate3D(-cam.X,-cam.Y,0)
	world := mgl32.Translate3D(
		item.Pos.X,
		item.Pos.Y,
		z,
	).Mul4(cxmath.Scale(worldItemSize))
	modelView := view.Mul4(world)
	spriteloader.DrawSpriteQuadMatrix(modelView,spriteId)
}

func (item *WorldItem) Tick(planet *world.Planet, dt float32) {
	item.Vel.Y -= physics.Gravity * dt
	item.Move(planet,dt)
}
