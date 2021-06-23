package item

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/models"
)

type ItemType struct {
	SpriteID int
	Name string

	Use func(ItemUseInfo)
}
type ItemTypeID uint32

type ItemUseInfo struct {
	Slot *InventorySlot
	ScreenX float32
	ScreenY float32
	Camera *camera.Camera
	Planet *world.Planet
	Player *models.Player
}

func (info ItemUseInfo) WorldCoords() mgl32.Vec2 {
	// click relative to camera
	camCoords :=
		mgl32.Vec4{
			info.ScreenX / render.PixelsPerTile,
			info.ScreenY / render.PixelsPerTile, 0, 1 }
	// click relative to world
	worldCoords := info.Camera.GetTransform().Mul4x1(camCoords)
	return worldCoords.Vec2()
}

func (info ItemUseInfo) PlayerCoords() mgl32.Vec2 {
	return mgl32.Vec2 { info.Player.Pos.X, info.Player.Pos.Y }
}

var itemTypes = []ItemType{}
func NewItemType(SpriteID int) ItemType {
	return ItemType {
		SpriteID: SpriteID,
		Name: "untitled",
		Use: func(ItemUseInfo) {},
	}
}

func AddItemType(itemType ItemType) ItemTypeID {
	id := ItemTypeID(len(itemTypes))
	itemTypes = append(itemTypes, itemType)
	return id
}

func GetItemTypeById(id ItemTypeID) *ItemType {
	return &itemTypes[id]
}
