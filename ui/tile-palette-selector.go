package ui

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/render"
)

type TilePaletteSelector struct {
	// store tiles for (1) displaying selector and (2) placing tiles
	Tiles             []world.Tile
	Transform         mgl32.Mat4
	Width             int
	SelectedTileIndex int
	visible           bool
	bluePixelSpriteId int
	redPixelSpriteId  int
}

const selectorSize = 2

func MakeTilePaleteSelector(tiles []world.Tile) TilePaletteSelector {
	width := math.Ceil(math.Sqrt(float64(len(tiles))))

	spriteloader.LoadSingleSprite("./assets/blue_pixel.png", "blue-pixel")
	spriteloader.LoadSingleSprite("./assets/red_pixel.png", "red-pixel")

	return TilePaletteSelector{
		Tiles:             tiles,
		Width:             int(width),
		SelectedTileIndex: -1,
		bluePixelSpriteId: spriteloader.GetSpriteIdByName("blue-pixel"),
		redPixelSpriteId:  spriteloader.GetSpriteIdByName("red-pixel"),
	}
}

const padding = 0.1
const spacing = 0.2

func (selector TilePaletteSelector) worldTransformForTileAtIndex(idx int) mgl32.Mat4 {
	offset := float32(selector.Width)/2 - 0.5
	yLocal := float32(idx/int(selector.Width)) - offset
	xLocal := float32(idx%int(selector.Width)) - offset
	tileScale := 1.0 / float32(selector.Width)
	localTransform :=
		mgl32.Scale3D(tileScale, tileScale, 1).
			Mul4(mgl32.Translate3D(xLocal, yLocal, 0))

	tileWorldTransform := selector.Transform.Mul4(localTransform)
	return tileWorldTransform
}

func (selector *TilePaletteSelector) UpdateTransform(ctx render.Context) {
	x := -ctx.Size.X()/2+float32(selector.Width)/2
	selector.Transform =
		mgl32.Translate3D(x,0,0).
		Mul4(cxmath.Scale(selectorSize))
}

func (selector *TilePaletteSelector) Draw(ctx render.Context) {
	selector.UpdateTransform(ctx)

	if !selector.visible {
		return
	}

	selectorCtx := ctx.PushLocal(selector.Transform)
	bgCtx := selectorCtx.PushLocal(cxmath.Scale(float32(1+padding)))

	spriteloader.
		DrawSpriteQuadContext(bgCtx, selector.redPixelSpriteId)

	if selector.SelectedTileIndex >= 0 {
		selectedTransform := selector.
			worldTransformForTileAtIndex(selector.SelectedTileIndex)

		spriteloader.
			DrawSpriteQuadMatrix(selectedTransform, selector.bluePixelSpriteId)
	}

	numTiles := float64(len(selector.Tiles))
	if numTiles > 0 {
		for idx, tile := range selector.Tiles {
			tileWorldTransform := selector.worldTransformForTileAtIndex(idx).
				Mul4(mgl32.Scale3D(1-spacing, 1-spacing, 1))

			// TODO add a custom texture for drawing air
			if tile.TileType != world.TileTypeNone {
				spriteloader.DrawSpriteQuadMatrix(tileWorldTransform, int(tile.SpriteID))
			}
		}
	}

}

func (selector *TilePaletteSelector) TrySelectTile(x, y float32) bool {
	// can't select palete tile when palete is invisible
	if !selector.visible {
		return false
	}

	// click relative to camera
	camCoords := mgl32.Vec4{x/render.PixelsPerTile,y/render.PixelsPerTile,0,1}
	// click relative to palete center in [-0.5,0.5]
	paleteCoords := selector.Transform.Inv().Mul4x1(camCoords).Vec2()
	// click relative to palete top left corner in [0,1]
	tileCoords := paleteCoords.Add(mgl32.Vec2{0.5,0.5})

	tileX := int(tileCoords.X()*float32(selector.Width))
	tileY := int(tileCoords.Y()*float32(selector.Width))

	if tileX >= 0 && tileX < selector.Width && tileY >= 0 && tileY < selector.Width {
		selector.SelectedTileIndex = tileY*selector.Width + tileX
		return true
	}
	return false
}

func (selector *TilePaletteSelector) GetSelectedTile() world.Tile {
	if selector.SelectedTileIndex >= 0 &&
		selector.SelectedTileIndex < len(selector.Tiles) {

		return selector.Tiles[selector.SelectedTileIndex]
	} else {
		return world.Tile{}
	}
}

func (selector *TilePaletteSelector) Toggle() {
	selector.visible = !selector.visible
}
