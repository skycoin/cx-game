package item

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
)

type Category uint32

const (
	Misc Category = iota
	BuildTool
)

type ItemType struct {
	SpriteID render.SpriteID
	Name     string
	Category Category

	Use func(ItemUseInfo)
	OnDrag func(ItemUseInfo, mgl32.Vec2, glfw.MouseButton)
	MouseDownRight func(ItemUseInfo) bool
}

//type ItemTypeID uint32
type ItemTypeID types.ItemTypeID

type ItemUseInfo struct {
	Slot      *InventorySlot
	ScreenX   float32
	ScreenY   float32
	Camera    *camera.Camera
	World     *world.World
	Player    *agents.Agent
	Inventory *Inventory
}

func (info ItemUseInfo) CamCoords() mgl32.Vec2 {
	return mgl32.Vec2{
		info.ScreenX / render.PixelsPerTile,
		info.ScreenY / render.PixelsPerTile,
	}
}

func (info ItemUseInfo) WorldCoords() mgl32.Vec2 {
	// click relative to camera
	camCoords := info.CamCoords().Vec4(0, 1)
	// click relative to world
	worldCoords := info.Camera.GetTransform().Mul4x1(camCoords)

	return info.World.Planet.WrapAround(worldCoords.Vec2())
}

func (info ItemUseInfo) PlayerCoords() mgl32.Vec2 {
	return mgl32.Vec2{
		info.Player.Transform.Pos.X, info.Player.Transform.Pos.Y}
}

var tileTypeIDsToItemTypeIDs = make(map[world.TileTypeID]ItemTypeID)
var itemTypeIDsToTileTypeIDs = make(map[ItemTypeID]world.TileTypeID)
var itemTypes = make(map[ItemTypeID]*ItemType)

var nextItemTypeID = ItemTypeID(1)

func NewItemType(SpriteID render.SpriteID) ItemType {
	return ItemType{
		SpriteID: SpriteID,
		Name:     "untitled-item",
		Use:      func(ItemUseInfo) {},
		OnDrag:     func(ItemUseInfo,mgl32.Vec2, glfw.MouseButton) {},
		MouseDownRight:      func(ItemUseInfo) bool {return false},
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

func GetTileTypeIDForItemTypeID(itemtypeID ItemTypeID) (world.TileTypeID, bool) {
	tiletypeID, ok := itemTypeIDsToTileTypeIDs[itemtypeID]
	return tiletypeID, ok
}

func GetItemTypeIdForTileTypeID(id world.TileTypeID) ItemTypeID {
	tiletype := id.Get()
	itemTypeID, ok := tileTypeIDsToItemTypeIDs[id]
	if ok {
		return itemTypeID
	}

	itemType := NewItemType((tiletype.ItemSpriteID))
	itemType.Name = tiletype.Name
	itemType.Use = func(info ItemUseInfo) {
		worldCoords := info.WorldCoords()
		x := int(worldCoords.X() + 0.5)
		y := int(worldCoords.Y() + 0.5)
		if !info.World.Planet.TileIsSolid(x, y) {
			info.Slot.Quantity--
			info.World.Planet.PlaceTileType(id, x, y)
		}
	}
	itemTypeID = AddItemType(itemType)
	tileTypeIDsToItemTypeIDs[id] = itemTypeID
	itemTypeIDsToTileTypeIDs[itemTypeID] = id
	return itemTypeID
}

func GetItemTypeIdForTile(tile world.Tile) ItemTypeID {
	return GetItemTypeIdForTileTypeID(tile.TileTypeID)
}

func (id ItemTypeID) Get() *ItemType {
	return itemTypes[id]
}
