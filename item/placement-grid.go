package item

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/cxmath/math32i"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
)

const PlacementGridWidth = 15
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
		for binIdx := range rects {
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
	canPlace              bool
}

func NewPlacementGrid() PlacementGrid {
	return PlacementGrid{PositionedTileTypeIDs: []PositionedTileTypeID{}}
}

func (grid *PlacementGrid) Assemble(toolType types.ToolType) {
	ids := world.TileTypeIDsForToolType(toolType)
	grid.PositionedTileTypeIDs = LayoutTiletypes(ids)
}

func (grid *PlacementGrid) Transform() mgl32.Mat4 {
	//slightly behind the hud to not overlap
	return mgl32.Translate3D(-10, grid.Scroll, constants.HUD_Z-1)
}

func (grid *PlacementGrid) Draw(ctx render.Context, invCam mgl32.Mat4) {
	ctx = ctx.PushLocal(grid.Transform())
	for _, positionedTileTypeID := range grid.PositionedTileTypeIDs {
		grid.DrawSlot(positionedTileTypeID, ctx)
	}
	if grid.HasSelected {
		previewColor := mgl32.Vec4{0, 1, 0, 0.5}
		if !grid.canPlace {
			previewColor = mgl32.Vec4{1, 0, 0, 0.5}
		}
		grid.DrawPreview(ctx, invCam, previewColor)
	}
}

// var previewColor = mgl32.Vec4{0, 1, 0, 0.5} // green
func (grid *PlacementGrid) DrawPreview(ctx render.Context, invCam mgl32.Mat4, previewColor mgl32.Vec4) {
	mousePos := input.GetMousePos().Mul(1.0 / float32(constants.PIXELS_PER_TILE))
	mousePosHomog := mgl32.Vec4{mousePos.X(), mousePos.Y(), 0, 1}
	mouseWorldPos := invCam.Inv().Mul4x1(mousePosHomog).Vec2()

	translate := mgl32.Translate3D(
		math32.Round(mouseWorldPos.X()),
		math32.Round(mouseWorldPos.Y()),
		0)
	shiftAndScale := grid.previewTransform()

	worldTransform := translate.Mul4(shiftAndScale)
	modelView := invCam.Mul4(worldTransform)

	//behind the hud
	modelView = modelView.Mul4(
		mgl32.Translate3D(0, 0, constants.HUD_Z-1),
	)
	render.DrawColorQuad(modelView, previewColor)
}

func (grid *PlacementGrid) previewTransform() mgl32.Mat4 {
	tiletype := grid.Selected.Get()

	unCenter :=
		mgl32.Translate3D(0.5, 0.5, 0)
	scaleUp :=
		mgl32.Scale3D(float32(tiletype.Width), float32(tiletype.Height), 1)
	reCenter :=
		mgl32.Translate3D(-0.5, -0.5, 0)

	return reCenter.Mul4(scaleUp).Mul4(unCenter)
}

func (grid PlacementGrid) DrawSlot(
	positionedTileTypeID PositionedTileTypeID, ctx render.Context,
) {
	slotCtx := ctx.PushLocal(positionedTileTypeID.Transform())
	// draw border
	color := borderColor
	if grid.Selected == positionedTileTypeID.TileTypeID {
		color = selectedBorderColor
	}
	render.DrawColorQuad(slotCtx.World, color)
	bgCtx := slotCtx.
		PushLocal(mgl32.Translate3D(0, 0, 0.1)).
		PushLocal(cxmath.Scale(1 - borderSize))
	render.DrawColorQuad(bgCtx.World, bgColor)
	// draw tiletype on top of bg
	itemCtx := slotCtx.PushLocal(cxmath.Scale(itemSize))

	render.DrawUISprite(
		itemCtx.World.Mul4(mgl32.Translate3D(0, 0, 0.2)),
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

func tilesAreClear(
	World *world.World, layerID world.LayerID,
	xstart, ystart, xstop, ystop int,
) bool {
	for x := xstart; x < xstop; x++ {
		for y := ystart; y < ystop; y++ {
			if !World.TileIsClear(layerID, x, y) {
				return false
			}
		}
	}
	// if layer is midlayer, do additional top layer check
	if layerID == world.MidLayer {
		for x := xstart; x < xstop; x++ {
			for y := ystart; y < ystop; y++ {
				if !World.TileIsClear(world.TopLayer, x, y) {
					return false
				}
			}
		}
	}
	return true
}

func (grid *PlacementGrid) TryPlace(info ItemUseInfo) bool {
	if !grid.HasSelected || grid.Selected == 0 {
		return false
	}
	worldCoords := info.WorldCoords()
	x32, y32 := cxmath.RoundVec2(worldCoords)
	x := int(x32)
	y := int(y32)
	if grid.canPlace {
		info.World.Planet.PlaceTileType(grid.Selected, x, y)
		return true
	}
	return false
}

func (grid *PlacementGrid) UpdatePreview(
		World *world.World, screenX, screenY float32, Cam *camera.Camera,
) {
	if !grid.HasSelected {
		return
	}

	mousePos := input.GetMousePos()
	screenTilePos := Cam.GetTransform().Mul4x1(mousePos.Mul(1.0/32).Vec4(0, 1))

	x32, y32 := cxmath.RoundVec2(screenTilePos.Vec2())
	x := int(x32)
	y := int(y32)

	tt := grid.Selected.Get()

	grid.canPlace = tilesAreClear(
		World,
		tt.Layer,
		x, y,
		x+int(tt.Width),
		y+int(tt.Height),
	)
}
