package item

import (
	"log"
	"strconv"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
)

const (
	inventoryScale float32 = 2.5
	hotbarScale    float32 = 2.5
	hotbarYOffset  float32 = 0.3
	// size of displayed item relative to slot
	itemSize   float32 = 0.8
	borderSize float32 = 0.1
)

type InventorySlot struct {
	ItemTypeID ItemTypeID
	Quantity   uint32
	Durability uint32
}

type Inventory struct {
	Width, Height        int
	Slots                []InventorySlot
	SelectedBarSlotIndex int
	GridHoldingIndex     int
	PlacementGrid        PlacementGrid

	IsOpen               bool
	lastMouseWorldPos    mgl32.Vec2
}

var (
	inventories         = []Inventory{}
	bgColor             = mgl32.Vec4{0.3, 0.3, 0.3, 1}
	borderColor         = mgl32.Vec4{0.8, 0.8, 0.8, 1}
	selectedBorderColor = mgl32.Vec4{0.8, 0.0, 0.0, 1}
)

func NewInventory(width, height int) types.InventoryID {
	inventories = append(inventories, Inventory{
		Width: width, Height: height,
		Slots:                make([]InventorySlot, width*height),
		SelectedBarSlotIndex: 3,
		PlacementGrid:        NewPlacementGrid(),
	})
	return types.InventoryID(len(inventories) - 1)
}

func GetInventoryById(id types.InventoryID) *Inventory {
	return &inventories[id]
}

//func (id types.InventoryID) Get() *Inventory { return &inventories[id] }

func (inventory Inventory) getBarSlots() []InventorySlot {
	return inventory.Slots[:inventory.Width]
}

func (inventory Inventory) getGridTransform() mgl32.Mat4 {
	s := inventoryScale
	return mgl32.Ident4().
		Mul4(mgl32.Translate3D(0, 0.5, 0)).
		Mul4(mgl32.Scale3D(s, s, s))
}

func (inventory Inventory) ItemTypeIDs() []ItemTypeID {
	ids := []ItemTypeID{}
	for _, slot := range inventory.Slots {
		if slot.Quantity > 0 {
			ids = append(ids, slot.ItemTypeID)
		}
	}
	return ids
}

func (inventory Inventory) DrawGrid(ctx render.Context) {
	gridTransform := inventory.getGridTransform()

	w := float32(inventory.Width)
	h := float32(inventory.Height)

	// draw all rows but first
	for y := 1; y < inventory.Height; y++ {
		for x := 0; x < inventory.Width; x++ {
			xRender := float32(x) - w/2
			yRender := h/2 - float32(y)
			slot := inventory.Slots[y*inventory.Width+x]
			slotTransform := gridTransform.
				Mul4(mgl32.Translate3D(xRender, yRender, 0))

			slotCtx := ctx.PushLocal(slotTransform)

			inventory.DrawSlot(slot, slotCtx, false)
		}
	}

	// draw last row a little further down
	y := 0
	for x := 0; x < inventory.Width; x++ {
		xRender := float32(x) - w/2
		yRender := h/2 - float32(inventory.Height) - 1
		slot := inventory.Slots[y*inventory.Width+x]
		slotTransform := gridTransform.
			Mul4(mgl32.Translate3D(xRender, yRender, 0))

		slotCtx := ctx.PushLocal(slotTransform)

		inventory.DrawSlot(slot, slotCtx, false)
	}
}

func (inv Inventory) DrawBar(ctx render.Context) {
	barCtx := ctx.
		PushLocal(mgl32.Translate3D(0, 1-ctx.Size.Y()/2, 0)).
		PushLocal(mgl32.Scale3D(hotbarScale, hotbarScale, 1))
	//barTransform := mgl32.Translate3D(0,-3,-spriteloader.SpriteRenderDistance)
	barSlots := inv.getBarSlots()
	for idx, slot := range barSlots {
		x := float32(idx) - float32(len(barSlots))/2
		slotCtx := barCtx.
			PushLocal(mgl32.Translate3D(x, hotbarYOffset, constants.HUD_Z))
		isSelected := idx == inv.SelectedBarSlotIndex
		inv.DrawSlot(slot, slotCtx, isSelected)
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
	transform := ctx.World.Mul4(mgl32.Translate3D(0,0,0.2))
	// draw border
	render.DrawColorQuad(transform, getBorderColor(isSelected))
	// draw bg on top of border
	bgTransform := transform.
		Mul4(mgl32.Translate3D(0, 0, 0.1)).
		Mul4(cxmath.Scale(1 - borderSize))
	render.DrawColorQuad(bgTransform, bgColor)
	// draw item on top of bg
	itemTransform := bgTransform.Mul4(cxmath.Scale(itemSize)).Mul4(
		mgl32.Translate3D(0,0,0.1))

	if slot.Quantity > 0 {
		spriteId := itemTypes[slot.ItemTypeID].SpriteID
		render.DrawUISprite(
			itemTransform.Mul4(
				mgl32.Translate3D(0, 0, 0.2)), spriteId,
			render.NewSpriteDrawOptions())

		textTransform := itemTransform.
			Mul4(mgl32.Translate3D(-0.5, 0.2, 0.3)).
			Mul4(cxmath.Scale(0.6))
		textCtx := render.Context {
			World: textTransform, Projection: ctx.Projection }
		if false {
			ui.DrawStringRightAligned(
				strconv.Itoa(int(slot.Quantity)),
				mgl32.Vec4{1, 1, 1, 1},
				textCtx,
			)
		}
		gl.Enable(gl.DEPTH_TEST)
		var x float32 = 320 + textCtx.MVP().At(0,3) * 320
		var y float32 = 240 - textCtx.MVP().At(1,3) * 240
		text := strconv.Itoa(int(slot.Quantity))
		ui.DefaultFont.Printf(x, y, 0.4, text)
	}
}

func (inventory Inventory) ItemSlotIndexForPosition(x, y int) int {
	return y*inventory.Width + x
}

func (inventory *Inventory) tryAddItemToStack(ItemTypeID ItemTypeID) bool {
	for idx, slot := range inventory.Slots {
		if slot.Quantity > 0 && slot.ItemTypeID == ItemTypeID {
			inventory.Slots[idx] = InventorySlot{
				ItemTypeID: ItemTypeID,
				Quantity:   slot.Quantity + 1,
			}
			return true
		}
	}
	return false
}

func (inventory *Inventory) tryAddItemToFreeSlot(ItemTypeID ItemTypeID) bool {
	for idx, slot := range inventory.Slots {
		if slot.Quantity == 0 {
			inventory.Slots[idx] = InventorySlot{
				ItemTypeID: ItemTypeID,
				Quantity:   1,
			}
			return true
		}
	}
	return false
}

func (inventory *Inventory) TryAddItem(ItemTypeID ItemTypeID) bool {
	return (inventory.tryAddItemToStack(ItemTypeID) ||
		inventory.tryAddItemToFreeSlot(ItemTypeID))
}

// X position of key along keyboard
func numberKeyPosition(k glfw.Key) int {
	if k == glfw.Key0 {
		return 9
	}
	return int(k-glfw.Key0) - 1
}

// Select a slot from the inventory bar based on the key pressed.
// slot layout matches keyboard layout.
// (1,2,3,4,5,6,7,8,9,0)
func (inventory *Inventory) TrySelectSlot(k glfw.Key) bool {
	if k < glfw.Key0 || k > glfw.Key9 {
		return false
	}

	selectedSlot := numberKeyPosition(k)
	if selectedSlot == inventory.SelectedBarSlotIndex {
		inventory.SelectedBarSlotIndex = -1
	} else {
		inventory.SelectedBarSlotIndex = selectedSlot
	}

	return true
}

func (inventory *Inventory) SelectedItemSlot() *InventorySlot {
	if inventory.SelectedBarSlotIndex == -1 {
		return nil
	}
	return &inventory.Slots[inventory.SelectedBarSlotIndex]
}

func (inventory *Inventory) TryUseItem(
	screenX, screenY float32, cam *camera.Camera,
	World *world.World, player *agents.Agent,
) bool {
	itemSlot := inventory.SelectedItemSlot()
	if itemSlot == nil { return false }
	// don't use empty items
	if itemSlot.Quantity == 0 { return false }
	itemType := GetItemTypeById(itemSlot.ItemTypeID)
	info := ItemUseInfo{
		Slot:    itemSlot,
		ScreenX: screenX, ScreenY: screenY,
		Camera:    cam,
		World:     World,
		Player:    player,
		Inventory: inventory,
	}
	itemType.Use(info)

	inventory.lastMouseWorldPos = info.WorldCoords()
	return true
}

func (inventory *Inventory) TryMouseDownRight(
	screenX, screenY float32, cam *camera.Camera,
	World *world.World, player *agents.Agent,
) bool {
	itemSlot := inventory.SelectedItemSlot()
	if itemSlot == nil { return false }
	// don't use empty items
	if itemSlot.Quantity == 0 { return false }
	itemType := GetItemTypeById(itemSlot.ItemTypeID)
	info := ItemUseInfo{
		Slot:    itemSlot,
		ScreenX: screenX, ScreenY: screenY,
		Camera:    cam,
		World:     World,
		Player:    player,
		Inventory: inventory,
	}
	consumed := itemType.MouseDownRight(info)
	if !consumed { return false }

	inventory.lastMouseWorldPos = info.WorldCoords()
	return true
}

func (inv *Inventory) TryDragItem(
	screenX, screenY float32, cam *camera.Camera,
	World *world.World, player *agents.Agent,
	b glfw.MouseButton,
) bool {
	itemSlot := inv.SelectedItemSlot()
	if itemSlot == nil { return false }
	// don't use empty items
	if itemSlot.Quantity == 0 { return false }
	itemType := GetItemTypeById(itemSlot.ItemTypeID)
	info := ItemUseInfo{
		Slot:    itemSlot,
		ScreenX: screenX, ScreenY: screenY,
		Camera:    cam,
		World:     World,
		Player:    player,
		Inventory: inv,
	}
	itemType.OnDrag(info, inv.lastMouseWorldPos, b)
	inv.lastMouseWorldPos = info.WorldCoords()
	return true
}

func (inventory *Inventory) getGridClickPosition(
	screenX, screenY float32,
) (idx int, ok bool) {
	camCoords := mgl32.Vec4{
		screenX / render.PixelsPerTile,
		screenY / render.PixelsPerTile,
		0, 1}

	centered := inventory.getGridTransform().Inv().Mul4x1(camCoords).Vec2()
	// convert from centered to top-left origin
	w := float32(inventory.Width)
	h := float32(len(inventory.Slots)) / w
	anchored := centered.Add(mgl32.Vec2{w / 2, h / 2})

	gridX := int(cxmath.Floor(anchored.X() + 0.5))
	gridY := int(cxmath.Floor(anchored.Y() + 0.5))

	if !inventory.IsOpen {
		if gridY == -1 {
			gridY = 0
		} else {
			return -1, false
		}
	}
	idx = -1
	clickIsOnGrid :=
		gridX >= 0 && gridX < inventory.Width &&
			// y=1 is a filler row - doesn't count
			gridY >= 0 && gridY <= int(h) // && gridY != 1

	if clickIsOnGrid {
		idx = inventory.SlotIdxForPosition(gridX, gridY)
	}

	return idx, clickIsOnGrid
}

func (inventory *Inventory) TryClickSlot(
	screenX, screenY float32, cam *camera.Camera,
	planet *world.Planet, player *agents.Agent,
) bool {
	idx, ok := inventory.getGridClickPosition(screenX, screenY)
	if ok {
		if !inventory.IsOpen {
			inventory.SelectedBarSlotIndex = idx
		} else {
			inventory.TrySelectGridSlot(idx)
		}
	}
	return ok
}

func (inventory *Inventory) OnReleaseMouse(
	screenX, screenY float32, cam *camera.Camera,
	planet *world.Planet, player *agents.Agent,
) bool {
	if !inventory.IsOpen {
		return false
	}
	idx, ok := inventory.getGridClickPosition(screenX, screenY)
	if ok {
		inventory.TryMoveGridSlot(idx)
	}

	return false
}

func (inventory *Inventory) TrySelectGridSlot(idx int) {
	slot := &inventory.Slots[idx]
	if slot.Quantity > 0 {
		inventory.GridHoldingIndex = idx
	}
}

func (inventory *Inventory) TryMoveGridSlot(idx int) {
	to := &inventory.Slots[idx]
	if to.Quantity == 0 && inventory.GridHoldingIndex >= 0 {
		from := &inventory.Slots[inventory.GridHoldingIndex]
		*to = *from
		from.Quantity = 0
	}
}

func (inventory *Inventory) SlotIdxForPosition(x, y int) int {
	// y=0 actually refers to last row - should probably fix this later
	if y == 0 {
		//y+=inventory.Height-1
	} else {
		y = inventory.Height - y // other rows are offset due to the gap
	}

	return y*inventory.Width + x
}

func (inv *Inventory) Draw(ctx render.Context, invCam mgl32.Mat4) {

	if inv.IsOpen {
		inv.DrawGrid(ctx)
	} else {
		inv.DrawBar(ctx)
	}
	slot := inv.SelectedItemSlot()
	if slot == nil {
		return
	}
	if slot.Quantity > 0 {
		item := slot.ItemTypeID.Get()
		if item.Category == BuildTool {
			// TODO do this less often
			toolType, ok := toolTypeFromItemName(item.Name)
			if !ok {
				log.Fatalf(
					"could not find tool type for string [%v]", item.Name)
			}
			inv.PlacementGrid.Assemble(toolType)
			inv.PlacementGrid.Draw(ctx, invCam)

		}
		// dev items
		if slot.ItemTypeID == EnemyToolItemTypeID {
			ui.DrawEnemyTool(ctx)
		}
	}
}

func (inv *Inventory) TryScrollDown() {
	inv.PlacementGrid.Scroll += placementGridScrollStride
}

func (inv *Inventory) TryScrollUp() {
	inv.PlacementGrid.Scroll -= placementGridScrollStride
}

func (inv *Inventory) TryCancelSelect() {
	inv.PlacementGrid.Selected = 0
}
