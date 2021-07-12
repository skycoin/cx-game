package item;

import (
	"log"
	"strconv"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/ui"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/models"
)

type InventorySlot struct {
	ItemTypeID ItemTypeID
	Quantity uint32
	Durability uint32
}

type Inventory struct {
	Width, Height int
	Slots []InventorySlot
	SelectedBarSlotIndex int
	GridHoldingIndex int
}

var inventories = []Inventory{}
var bgColor = mgl32.Vec4{0.3,0.3,0.3,1}
var borderColor = mgl32.Vec4{0.8,0.8,0.8,1}
var selectedBorderColor = mgl32.Vec4{0.8,0,0,1}

func NewInventory(width, height int) uint32 {
	inventories = append(inventories, Inventory {
		Width: width, Height: height,
		Slots: make([]InventorySlot, width*height),
		SelectedBarSlotIndex: 3,
	})
	return uint32(len(inventories)-1)
}

func NewDevInventory() uint32 {
	inventoryId := NewInventory(10, 8)
	inventory := GetInventoryById(inventoryId)
	inventory.Slots[inventory.ItemSlotIndexForPosition(1, 0)] =
		InventorySlot{LaserGunItemTypeID, 1, 0}
	inventory.Slots[inventory.ItemSlotIndexForPosition(2, 0)] =
		InventorySlot{GunItemTypeID, 1, 0}

	pipeTileType,ok := world.GetTileTypeByID(world.TileTypeIDs.Pipe)
	if !ok { log.Fatal("Cannot find pipe tile type")}
	pipeTile := pipeTileType.CreateTile(world.TileCreationOptions{})
	pipeItemTypeID := GetItemTypeIdForTile(pipeTile)
	inventory.Slots[inventory.ItemSlotIndexForPosition(3, 0)] =
		InventorySlot{pipeItemTypeID, 20, 0}

	return inventoryId
}

func GetInventoryById(id uint32) *Inventory {
	return &inventories[id]
}

func (inventory Inventory) getBarSlots() []InventorySlot {
	return inventory.Slots[:inventory.Width]
}

func (inventory Inventory) getGridTransform() mgl32.Mat4 {
	return mgl32.Ident4().
		Mul4(mgl32.Translate3D(0,0.5,0)).
		Mul4(mgl32.Scale3D(gridScale,gridScale,gridScale))
}

var gridScale float32 = 1.5
// size of displayed item relative to slot
var itemSize float32 = 0.8
var borderSize float32 = 0.1
func (inventory Inventory) DrawGrid(ctx render.Context) {
	gridTransform := inventory.getGridTransform()

	w := float32(inventory.Width)
	h := float32(inventory.Height)

	// draw all rows but first
	for y:=1;y<inventory.Height;y++ {
		for x:=0;x<inventory.Width;x++ {
			xRender := float32(x) - w / 2
			yRender := h / 2 - float32(y)
			slot := inventory.Slots[y*inventory.Width+x]
			slotTransform := gridTransform.
				Mul4(mgl32.Translate3D(xRender,yRender,0))

			slotCtx := ctx.PushLocal(slotTransform)

			inventory.DrawSlot(slot,slotCtx,false)
		}
	}

	// draw last row a little further down
	y := 0
	for x:=0;x<inventory.Width;x++ {
		xRender := float32(x) - w / 2
		yRender := h / 2 - float32(inventory.Height) - 1
		slot := inventory.Slots[y*inventory.Width+x]
		slotTransform := gridTransform.
			Mul4(mgl32.Translate3D(xRender,yRender,0))

		slotCtx := ctx.PushLocal(slotTransform)

		inventory.DrawSlot(slot,slotCtx,false)
	}
}

func (inv Inventory) DrawBar(ctx render.Context) {
	barCtx := ctx.PushLocal(mgl32.Translate3D(0,1-ctx.Size.Y()/2,0))
	//barTransform := mgl32.Translate3D(0,-3,-spriteloader.SpriteRenderDistance)
	barSlots := inv.getBarSlots()
	for idx,slot := range barSlots {
		x := float32(idx) - float32(len(barSlots)) / 2
		slotCtx := barCtx.PushLocal(mgl32.Translate3D(x,0,0))
		isSelected := idx == inv.SelectedBarSlotIndex
		inv.DrawSlot(slot,slotCtx,isSelected)
	}
}

func getBorderColor(isSelected bool) mgl32.Vec4 {
	if isSelected {
		return selectedBorderColor
	} else {
		return borderColor
	}
}

func (inventory Inventory) DrawSlot(
		slot InventorySlot, ctx render.Context, isSelected bool,
) {
	// draw border
	utility.DrawColorQuad(ctx,getBorderColor(isSelected))
	// draw bg on top of border
	bgCtx := ctx.PushLocal(cxmath.Scale(1-borderSize))
	utility.DrawColorQuad(bgCtx,bgColor)
	// draw item on top of bg
	itemCtx := ctx.PushLocal(cxmath.Scale(itemSize))
	// TODO write number for quantity
	if slot.Quantity > 0 {
		spriteId := itemTypes[slot.ItemTypeID].SpriteID
		spriteloader.DrawSpriteQuadContext(
			itemCtx, (spriteId), spriteloader.NewDrawOptions() )

		textCtx := itemCtx.PushLocal(
			mgl32.Translate3D(0.5,-0.05,0).
			Mul4(cxmath.Scale(0.6)))
		ui.DrawStringRightAligned(
			strconv.Itoa(int(slot.Quantity)),
			mgl32.Vec4 {1,1,1,1},
			textCtx,
		)
	}
}

func (inventory Inventory) ItemSlotIndexForPosition(x,y int) int {
	return y*inventory.Width+x
}

func (inventory *Inventory) tryAddItemToStack(ItemTypeID ItemTypeID) bool {
	for idx,slot := range inventory.Slots {
		if slot.Quantity>0 && slot.ItemTypeID == ItemTypeID {
			inventory.Slots[idx] = InventorySlot {
				ItemTypeID: ItemTypeID,
				Quantity: slot.Quantity+1,
			}
			return true
		}
	}
	return false
}

func (inventory *Inventory) tryAddItemToFreeSlot(ItemTypeID ItemTypeID) bool {
	for idx,slot := range inventory.Slots {
		if slot.Quantity == 0 {
			inventory.Slots[idx] = InventorySlot {
				ItemTypeID: ItemTypeID,
				Quantity: 1,
			}
			return true
		}
	}
	return false
}

func (inventory *Inventory) TryAddItem(ItemTypeID ItemTypeID) bool {
	return (
		inventory.tryAddItemToStack(ItemTypeID) ||
		inventory.tryAddItemToFreeSlot(ItemTypeID) )
}

// Select a slot from the inventory bar based on the key pressed.
// slot layout matches keyboard layout.
// (1,2,3,4,5,6,7,8,9,0)
func (inventory *Inventory) TrySelectSlot(k glfw.Key) bool {
	if k < glfw.Key0 || k > glfw.Key9 {
		return false
	}

	idx := int(k-glfw.Key0)-1
	if idx == -1 {
		idx = 9
	}
	inventory.SelectedBarSlotIndex = idx

	return true
}

func (inventory *Inventory) SelectedItemSlot() *InventorySlot {
	return &inventory.Slots[inventory.SelectedBarSlotIndex]
}

func (inventory *Inventory) TryUseItem(
		screenX,screenY float32, cam *camera.Camera,
		planet *world.Planet, player *models.Player,
) bool {
	itemSlot := inventory.SelectedItemSlot()
	// don't use empty items
	if itemSlot.Quantity == 0 {
		return false
	}
	itemType := GetItemTypeById(itemSlot.ItemTypeID)
	itemType.Use(ItemUseInfo {
		Slot: itemSlot,
		ScreenX: screenX, ScreenY: screenY,
		Camera: cam, Planet: planet,
		Player: player,
	})
	return true
}

func (inventory *Inventory) getGridClickPosition(
		screenX,screenY float32,
) (idx int, ok bool) {
	camCoords := mgl32.Vec4 {
		screenX / render.PixelsPerTile,
		screenY / render.PixelsPerTile,
		0, 1 }

	centered := inventory.getGridTransform().Inv().Mul4x1(camCoords).Vec2()
	// convert from centered to top-left origin
	w := float32(inventory.Width)
	h := float32(len(inventory.Slots)) / w
	anchored := centered.Add(mgl32.Vec2 { w/2, h/2 })

	gridX := int(anchored.X()+0.5)
	gridY := int(anchored.Y()+0.5)

	idx = -1
	clickIsOnGrid :=
		gridX >=0 && gridX < inventory.Width &&
		// y=1 is a filler row - doesn't count
		gridY >=0 && gridY <= int(h) && gridY != 1

	if clickIsOnGrid {
		idx = inventory.SlotIdxForPosition(gridX,gridY)
	}

	return idx,clickIsOnGrid
}

func (inventory *Inventory) TryClickSlot(
		screenX,screenY float32, cam *camera.Camera,
		planet *world.Planet, player *models.Player,
) bool {
	idx,ok := inventory.getGridClickPosition(screenX,screenY)
	if ok {
		inventory.TrySelectGridSlot(idx)
	}
	
	return ok
}

func (inventory *Inventory) TryMoveSlot(
		screenX,screenY float32, cam *camera.Camera,
		planet *world.Planet, player *models.Player,
) bool {
	idx,ok := inventory.getGridClickPosition(screenX,screenY)
	if ok {
		inventory.TryMoveGridSlot(idx)
	}
	
	return false
}

func (inventory *Inventory) TrySelectGridSlot(idx int) {
	slot := &inventory.Slots[idx]
	if slot.Quantity>0 {
		inventory.GridHoldingIndex = idx
	}
}

func (inventory *Inventory) TryMoveGridSlot(idx int) {
	to := &inventory.Slots[idx]
	if to.Quantity==0 && inventory.GridHoldingIndex>=0 {
		from := &inventory.Slots[inventory.GridHoldingIndex]
		*to = *from
		from.Quantity=0
	}
}

func (inventory *Inventory) SlotIdxForPosition(x,y int) int {
	// y=0 actually refers to last row - should probably fix this later
	if y==0 {
		//y+=inventory.Height-1
	} else {
		y=inventory.Height-y // other rows are offset due to the gap
	}

	return y*inventory.Width + x
}
