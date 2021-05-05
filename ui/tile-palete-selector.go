package ui

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/spriteloader"
	//"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/cxmath"
)

type TilePaleteSelector struct {
	// store tiles for (1) displaying selector and (2) placing tiles
	Tiles []world.Tile
	Transform mgl32.Mat4
	Width int
	SelectedTileIndex int
	visible bool
	bluePixelSpriteId int
	redPixelSpriteId int
}

func MakeTilePaleteSelector(tiles []world.Tile) TilePaleteSelector {
	width := math.Ceil(math.Sqrt(float64(len(tiles))))

	spriteloader.LoadSingleSprite("./assets/blue_pixel.png","blue-pixel")
	spriteloader.LoadSingleSprite("./assets/red_pixel.png","red-pixel")

	return TilePaleteSelector {
		Tiles: tiles,
		Transform: mgl32.Translate3D(-5,0,-spriteloader.SpriteRenderDistance),
		Width: int(width),
		SelectedTileIndex: -1,
		bluePixelSpriteId: spriteloader.GetSpriteIdByName("blue-pixel"),
		redPixelSpriteId: spriteloader.GetSpriteIdByName("red-pixel"),
	}
}

const padding = 0.1
const spacing = 0.2

func (selector TilePaleteSelector) worldTransformForTileAtIndex(idx int) mgl32.Mat4 {
	offset := float32(selector.Width)/2 - 0.5
	yLocal := float32(idx/int(selector.Width))-offset
	xLocal := float32(idx%int(selector.Width))-offset
	tileScale := 1.0/float32(selector.Width)
	localTransform :=
		mgl32.Scale3D(tileScale,tileScale,1).
		Mul4(mgl32.Translate3D(xLocal,yLocal,0))

	tileWorldTransform := selector.Transform.Mul4(localTransform)
	return tileWorldTransform
}

func (selector *TilePaleteSelector) Draw() {
	if !selector.visible {
		return
	}

	// TODO create a shader for drawing arbitrary colors
	bgScale := float32(1+padding)
	bgTransform := mgl32.Mat4.Mul4(
		selector.Transform,
		mgl32.Scale3D(bgScale,bgScale,1),
	)

	spriteloader.
		DrawSpriteQuadMatrix(bgTransform,selector.redPixelSpriteId)

	if selector.SelectedTileIndex >=0 {
		selectedTransform := selector.
			worldTransformForTileAtIndex(selector.SelectedTileIndex)

		spriteloader.
			DrawSpriteQuadMatrix(selectedTransform,selector.bluePixelSpriteId)
	}

	numTiles := float64(len(selector.Tiles))
	if numTiles>0 {
		for idx,tile := range selector.Tiles {
			tileWorldTransform := selector.worldTransformForTileAtIndex(idx).
				Mul4(mgl32.Scale3D(1-spacing,1-spacing,1))

			// TODO add a custom texture for drawing air
			if tile.TileType!=world.TileTypeNone {
				spriteloader.DrawSpriteQuadMatrix(tileWorldTransform,int(tile.SpriteID))
			}
		}
	}

}

func (selector *TilePaleteSelector) TrySelectTile(x,y float32, projection mgl32.Mat4) bool {
	// can't select palete tile when palete is invisible
	if !selector.visible {
		return false
	}
	worldCoords := cxmath.ConvertScreenCoordsToWorld(x,y,projection)
	paleteCoords := selector.Transform.Inv().Mul4x1(worldCoords).Vec2()
	tileX := int(math.Floor(float64(paleteCoords.X()+0.5)*float64(selector.Width)))
	tileY := int(math.Floor(float64(paleteCoords.Y()+0.5)*float64(selector.Width)))

	if tileX>=0 && tileX<selector.Width && tileY>=0 && tileY<selector.Width {
		selector.SelectedTileIndex = tileY*selector.Width + tileX
		return true
	}
	return false
}

func (selector *TilePaleteSelector) GetSelectedTile() world.Tile {
	if selector.SelectedTileIndex>=0 &&
		selector.SelectedTileIndex<len(selector.Tiles) {

		return selector.Tiles[selector.SelectedTileIndex]
	} else {
		return world.Tile {}
	}
}


func (selector *TilePaleteSelector) Toggle() {
	selector.visible = !selector.visible
}
