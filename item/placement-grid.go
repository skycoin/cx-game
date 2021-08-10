package item

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32i"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
)

const PlacementGridWidth = 5
const placementGridScrollStride = 4

func binTileTypeIDsByMaterial(
	tiletypeIDs []world.TileTypeID,
) map[world.MaterialID][]world.TileTypeID {
	bins := make(map[world.MaterialID][]world.TileTypeID)
	for _, tiletypeID := range tiletypeIDs {
		_, ok := bins[tiletypeID.Get().MaterialID]
		if !ok {
			bins[tiletypeID.Get().MaterialID] = []world.TileTypeID{}
		}
		bins[tiletypeID.Get().MaterialID] =
			append(bins[tiletypeID.Get().MaterialID], tiletypeID)
	}
	return bins
}

func GetTileTypesIDsForItemTypeIDs(
	itemtypeIDs []ItemTypeID,
) []world.TileTypeID {
	tiletypeIDs := []world.TileTypeID{}
	for _, itemtypeID := range itemtypeIDs {
		tiletypeID, ok := GetTileTypeIDForItemTypeID(itemtypeID)
		if ok {
			tiletypeIDs = append(tiletypeIDs, tiletypeID)
		}
	}
	return tiletypeIDs
}

type PositionedTileTypeID struct {
	TileTypeID world.TileTypeID
	Rect       cxmath.Rect
}

func (p PositionedTileTypeID) Transform() mgl32.Mat4 {
	translation := mgl32.Translate3D(
		float32(p.Rect.Size.X)/2-0.5+float32(p.Rect.Origin.X),
		-1*(float32(p.Rect.Size.Y)/2-0.5+float32(p.Rect.Origin.Y)),
		0)
	scale := mgl32.Scale3D(
		float32(p.Rect.Size.X), float32(p.Rect.Size.Y), 1)
	return translation.Mul4(scale)
}

func getTileTypeSizes(ids []world.TileTypeID) []cxmath.Vec2i {
	sizes := make([]cxmath.Vec2i, len(ids))
	for idx, id := range ids {
		sizes[idx] = id.Get().Size()
	}
	return sizes
}

func LayoutTiletypes(tiletypeIDs []world.TileTypeID) []PositionedTileTypeID {
	bins := binTileTypeIDsByMaterial(tiletypeIDs)
	positionedTileTypeIDs := make([]PositionedTileTypeID, len(tiletypeIDs))
	layoutIdx := 0

	materialYOffset := int32(0)

	for _, bin := range bins {
		sizes := getTileTypeSizes(bin)
		rects := cxmath.PackRectangles(PlacementGridWidth, sizes)
		for binIdx, _ := range rects {
			rect := &rects[binIdx]
			rect.Origin.Y += materialYOffset
			positionedTileTypeIDs[layoutIdx] = PositionedTileTypeID{
				TileTypeID: tiletypeIDs[layoutIdx],
				Rect:       *rect,
			}
			layoutIdx++
		}
		for _, rect := range rects {
			materialYOffset = math32i.Max(materialYOffset, rect.Bottom())
		}
	}

	return positionedTileTypeIDs
}

type PlacementGrid struct {
	PositionedTileTypeIDs []PositionedTileTypeID
	Selected              world.TileTypeID
	HasSelected           bool
	Scroll                float32
}

func NewPlacementGrid() PlacementGrid {
	return PlacementGrid{PositionedTileTypeIDs: []PositionedTileTypeID{}}
}

func (grid *PlacementGrid) Assemble(itemTypeIDs []ItemTypeID) {
	ids := world.AllTileTypeIDs()
	grid.PositionedTileTypeIDs = LayoutTiletypes(ids)
}

func (grid *PlacementGrid) Transform() mgl32.Mat4 {
	return mgl32.Translate3D(-10, grid.Scroll, 0)
}

func (ig *PlacementGrid) Draw(ctx render.Context) {
	ctx = ctx.PushLocal(ig.Transform())
	for _, positionedTileTypeID := range ig.PositionedTileTypeIDs {
		ig.DrawSlot(positionedTileTypeID, ctx)
	}
}

func (ig PlacementGrid) DrawSlot(
	positionedTileTypeID PositionedTileTypeID, ctx render.Context,
) {
	slotCtx := ctx.PushLocal(positionedTileTypeID.Transform())
	// draw border
	render.DrawColorQuad(slotCtx, borderColor)
	bgCtx := slotCtx.
		PushLocal(mgl32.Translate3D(0,0,0.1)).
		PushLocal(cxmath.Scale(1 - borderSize))
	render.DrawColorQuad(bgCtx, bgColor)
	// draw tiletype on top of bg
	itemCtx := slotCtx.PushLocal(cxmath.Scale(itemSize))

	render.DrawUISprite(
		itemCtx.World.Mul4(mgl32.Translate3D(0,0,0.2)),
		positionedTileTypeID.TileTypeID.Get().ItemSpriteID,
		render.NewSpriteDrawOptions(),
	)
}

func (grid *PlacementGrid) TrySelect(camCoords mgl32.Vec2) bool {
	relative := grid.Transform().Inv().Mul4x1(camCoords.Vec4(0, 1)).Vec2()
	x, y := cxmath.RoundVec2(relative)
	y = -y
	for _, positioned := range grid.PositionedTileTypeIDs {
		if positioned.Rect.Contains(x, y) {
			grid.Selected = positioned.TileTypeID
			grid.HasSelected = true
			return true
		}
	}
	return false
}

func (grid *PlacementGrid) TryPlace(info ItemUseInfo) bool {
	if !grid.HasSelected {
		return false
	}
	worldCoords := info.WorldCoords()
	x32, y32 := cxmath.RoundVec2(worldCoords)
	x := int(x32)
	y := int(y32)
	if info.World.TileIsClear(x, y) {
		info.World.Planet.PlaceTileType(grid.Selected, x, y)
		return true
	}
	return false
}
