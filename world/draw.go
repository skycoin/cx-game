package world

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/mathi"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/render"
)

type PositionedTile struct {
	Tile     Tile
	Position cxmath.Vec2i
}

func (pt PositionedTile) Transform() mgl32.Mat4 {
	translate := mgl32.Translate3D(
		float32(pt.Position.X), float32(pt.Position.Y), 0,
	)
	tiletype := pt.Tile.TileTypeID.Get()
	return translate.Mul4(tiletype.Transform())
}

func (planet *Planet) DrawLayer(
	layer Layer, cam *camera.Camera, layerID LayerID,
) {
	planet.program.Use()
	defer planet.program.StopUsing()

	w := int(planet.Width)
	// split up planet into 2 hemispheres to achieve wrap around
	// without calculating relative tile positions individually
	planet.DrawHemisphere(layer, cam, 0, w/2, layerID)
	planet.DrawHemisphere(layer, cam, w/2, w, layerID)
}

func (planet *Planet) DrawHemisphere(
	layer Layer, cam *camera.Camera, left, right int, layerID LayerID,
) {
	center := float32((left + right) / 2)
	_ = center

	projection := render.Projection
	planet.program.Use()
	planet.program.
		SetMat4("projection", &projection)
	planet.program.StopUsing()
	planet.liquidProgram.Use()
	planet.liquidProgram.
		SetMat4("projection", &projection)
	planet.liquidProgram.StopUsing()
	planet.program.Use()

	visible := planet.visibleTiles(layer, cam, left, right)
	for _, positionedTile := range visible {
		z := layerID.Z()
		transform := positionedTile.Transform().
			Mul4(mgl32.Translate3D(0, 0, z)).
			Mul4(positionedTile.Tile.FlipTransform)
		drawOpts := render.NewSpriteDrawOptions()
		if layerID == TopLayer {
			drawOpts.Outline = true
		}
		if layerID == WindowLayer {
			drawOpts.Cutout = true
		}
		if layerID == SuperLayer {
			drawOpts.Translucent = true
		}
		render.
			DrawWorldSprite(transform, positionedTile.Tile.SpriteID, drawOpts)
	}
}

func filterLiquidTiles(all []PositionedTile) []PositionedTile {
	liquids := []PositionedTile{}
	for _, tile := range all {
		if tile.Tile.TileCategory == TileCategoryLiquid {
			liquids = append(liquids, tile)
		}
	}
	return liquids
}

func (planet *Planet) Draw(cam *camera.Camera, layerID LayerID) {
	planet.DrawLayer(planet.Layers[layerID], cam, layerID)
}

func (planet *Planet) visibleTiles(
	layer Layer, cam *camera.Camera, left, right int,
) []PositionedTile {
	bottom := mathi.Max(cam.Frustum.Bottom, 0)
	top := mathi.Min(cam.Frustum.Top, int(planet.Height))
	capacity := (top - bottom + 1) * (right - left + 1)
	if capacity < 0 {
		capacity = 0
	}
	positionedTiles := make([]PositionedTile, 0, capacity)

	for y := bottom; y <= top; y++ {
		for x := left; x <= right; x++ {
			tileIdx := planet.GetTileIndex(x, y)
			if tileIdx != -1 {
				tile := layer.Tiles[tileIdx]
				shouldRender :=
					tile.TileCategory != TileCategoryNone &&
						tile.TileCategory != TileCategoryChild
				if shouldRender {
					positionedTiles = append(positionedTiles, PositionedTile{
						Position: cxmath.Vec2i{X: int32(x), Y: int32(y)},
						Tile:     tile,
					})
				}
			}
		}
	}
	return positionedTiles
}
