package item

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/ids"
)

type ItemType struct {
	SpriteID spriteloader.SpriteID
	Name string

	Use func(ItemUseInfo)
}
//type ItemTypeID uint32
type ItemTypeID ids.ItemTypeID

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

	return info.Planet.WrapAround(worldCoords.Vec2())
}

func (info ItemUseInfo) PlayerCoords() mgl32.Vec2 {
	return mgl32.Vec2 { info.Player.Pos.X, info.Player.Pos.Y }
}

var tileTypeIDsToItemTypeIDs = make(map[world.TileTypeID]ItemTypeID)
var itemTypes = make(map[ItemTypeID]*ItemType)

var nextItemTypeID = ItemTypeID(1)

func NewItemType(SpriteID spriteloader.SpriteID) ItemType {
	return ItemType {
		SpriteID: SpriteID,
		Name: "untitled",
		Use: func(ItemUseInfo) {},
	}
}

func AddItemType(itemType ItemType) ItemTypeID {
	id := nextItemTypeID
	itemTypes[id] = &itemType
	nextItemTypeID++
	return id
}

func GetItemTypeById(id ItemTypeID) *ItemType {
	return itemTypes[id]
}

func GetItemTypeIdForTileTypeID(id world.TileTypeID) ItemTypeID {
	tiletype := id.Get()
	itemTypeID,ok := tileTypeIDsToItemTypeIDs[id]
	if ok { return itemTypeID }

	itemType := NewItemType((tiletype.ItemSpriteID))
	itemType.Name = tiletype.Name
	itemType.Use = func(info ItemUseInfo) {
		worldCoords := info.WorldCoords()
		x := int(worldCoords.X()+0.5)
		y := int(worldCoords.Y()+0.5)
		if !info.Planet.TileIsSolid(x,y) {
			info.Slot.Quantity--
			info.Planet.PlaceTileType(id,x,y)
		}
	}
	itemTypeID = AddItemType(itemType)
	tileTypeIDsToItemTypeIDs[id] = itemTypeID
	return itemTypeID
}

func GetItemTypeIdForTile(tile world.Tile) ItemTypeID {
	return GetItemTypeIdForTileTypeID(tile.TileTypeID)
}
