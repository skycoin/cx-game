package item;

import (
	"strconv"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/ui"
)

type InventorySlot struct {
	ItemTypeID uint32
	Quantity uint32
}

type Inventory struct {
	Width, Height int
	Slots []InventorySlot
}

var inventories = []Inventory{}
var bgColor = mgl32.Vec4{0.3,0.3,0.3,1}
var borderColor = mgl32.Vec4{0.8,0.8,0.8,1}

func NewInventory(width, height int) uint32 {
	inventories = append(inventories, Inventory {
		Width: width, Height: height,
		Slots: make([]InventorySlot, width*height),
	})
	return uint32(len(inventories)-1)
}

func GetInventoryById(id uint32) *Inventory {
	return &inventories[id]
}

// TODO really need to figure out screen space solution

func (inventory Inventory) getBarSlots() []InventorySlot {
	start := inventory.Width*(inventory.Height-1)
	return inventory.Slots[start:]
}

var gridScale float32 = 0.8
// size of displayed item relative to slot
var itemSize float32 = 0.8
var borderSize float32 = 0.1
func (inventory Inventory) DrawGrid() {
	gridTransform :=
		mgl32.Translate3D(0,0.5,-spriteloader.SpriteRenderDistance).
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

			inventory.DrawSlot(slot,slotTransform)
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

		inventory.DrawSlot(slot,slotTransform)
	}
}

func (inv Inventory) DrawBar() {
	barTransform := mgl32.Translate3D(0,-3,-spriteloader.SpriteRenderDistance)
	barSlots := inv.getBarSlots()
	for idx,slot := range barSlots {
		x := float32(idx) - float32(len(barSlots)) / 2
		slotTransform := barTransform.
			Mul4(mgl32.Translate3D(x,0,0))
		_ = idx
		_ = slot
		// TODO draw the correct sprite
		inv.DrawSlot(slot,slotTransform)
	}
}

func (inventory Inventory) DrawSlot(slot InventorySlot, world mgl32.Mat4) {
	// draw border
	utility.DrawColorQuad(world,borderColor)
	// draw bg on top of border
	bgTransform := world.Mul4(cxmath.Scale(1-borderSize))
	utility.DrawColorQuad(bgTransform,bgColor)
	// draw item on top of bg
	itemTransform := world.Mul4(cxmath.Scale(itemSize))
	// TODO write number for quantity
	if slot.Quantity > 0 {
		spriteId := itemTypes[slot.ItemTypeID].SpriteID
		spriteloader.DrawSpriteQuadMatrix(itemTransform, int(spriteId))

		textTransform := itemTransform.
			Mul4(mgl32.Translate3D(0.5,-0.05,0)).
			Mul4(cxmath.Scale(0.6))
		ui.DrawStringRightAligned(
			strconv.Itoa(int(slot.Quantity)),
			textTransform,
			mgl32.Vec4 {1,1,1,1},
		)
	}
}

func (inventory Inventory) ItemSlotIndexForPosition(x,y int) int {
	return y*inventory.Width+x
}

func (inventory *Inventory) tryAddItemToStack(ItemTypeID uint32) bool {
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

func (inventory *Inventory) tryAddItemToFreeSlot(ItemTypeID uint32) bool {
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

func (inventory *Inventory) TryAddItem(ItemTypeID uint32) bool {
	return (
		inventory.tryAddItemToStack(ItemTypeID) ||
		inventory.tryAddItemToFreeSlot(ItemTypeID) )
}
