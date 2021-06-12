package item

import (
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
