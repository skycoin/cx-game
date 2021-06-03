package item;

import (
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

func GetInventoryById(id uint32) *Inventory {
	return &inventories[id]
}

func (inventory Inventory) getBarSlots() []InventorySlot {
	start := inventory.Width*(inventory.Height-1)
	return inventory.Slots[start:]
}

var gridScale float32 = 1.5
// size of displayed item relative to slot
var itemSize float32 = 0.8
var borderSize float32 = 0.1
func (inventory Inventory) DrawGrid(ctx render.Context) {
	gridTransform :=
		mgl32.Translate3D(0,0.5,0).
		Mul4(mgl32.Scale3D(gridScale,gridScale,gridScale))

	w := float32(inventory.Width)
	h := float32(inventory.Height)

	// draw all rows but last
	for y:=0;y<inventory.Height-1;y++ {
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
	y := inventory.Height-1
	for x:=0;x<inventory.Width;x++ {
		xRender := float32(x) - w / 2
		yRender := h / 2 - float32(y) - 1
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
		spriteloader.DrawSpriteQuadContext(itemCtx, int(spriteId))

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
	globalIdx :=
		inventory.SelectedBarSlotIndex +
		inventory.Width*(inventory.Height-1)
	return &inventory.Slots[globalIdx]
}

func (inventory *Inventory) TryUseItem(
		screenX,screenY float32, cam *camera.Camera,
		planet *world.Planet, player *models.Cat,
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
