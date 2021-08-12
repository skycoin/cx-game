package item

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32i"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/engine/input"
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

func (grid *PlacementGrid) Draw(ctx render.Context, camPos mgl32.Vec2) {
	ctx = ctx.PushLocal(grid.Transform())
	for _, positionedTileTypeID := range grid.PositionedTileTypeIDs {
		grid.DrawSlot(positionedTileTypeID, ctx)
	}
	if grid.HasSelected { grid.DrawPreview(ctx, camPos) }
}

var previewColor = mgl32.Vec4 { 0,1,0,0.5 } // green
func (grid *PlacementGrid) DrawPreview(ctx render.Context, camPos mgl32.Vec2) {
	mousePos := input.GetMousePos().
		Mul(1.0/float32(constants.PIXELS_PER_TILE)).
		Add(mgl32.Vec2{0.5,0.5})

	mouseWorldPos := camPos.Add(mousePos)
	offsetIntoTile := mgl32.Vec2 {
		math32.Mod( mouseWorldPos.X(), 1),
		math32.Mod( mouseWorldPos.Y(), 1),
	}
	tilePos := mousePos.Sub(offsetIntoTile)

	translate := mgl32.Translate3D(tilePos.X(), tilePos.Y(), 0)
	scaleAndShift := grid.previewTransform()
	modelView := translate.Mul4(scaleAndShift)

	render.DrawColorQuad(modelView, previewColor)
}

func (grid *PlacementGrid) previewTransform() mgl32.Mat4 {
	tiletype := grid.Selected.Get()

	unCenter :=
		mgl32.Translate3D( 0.5, 0.5, 0)
	scaleUp :=
		mgl32.Scale3D( float32(tiletype.Width), float32(tiletype.Height), 1, )
	reCenter :=
		mgl32.Translate3D( -0.5, -0.5, 0)

	return reCenter.Mul4(scaleUp).Mul4(unCenter)
}

func (ig PlacementGrid) DrawSlot(
	positionedTileTypeID PositionedTileTypeID, ctx render.Context,
) {
	slotCtx := ctx.PushLocal(positionedTileTypeID.Transform())
	// draw border
	render.DrawColorQuad(slotCtx.World, borderColor)
	bgCtx := slotCtx.
		PushLocal(mgl32.Translate3D(0,0,0.1)).
		PushLocal(cxmath.Scale(1 - borderSize))
	render.DrawColorQuad(bgCtx.World, bgColor)
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
