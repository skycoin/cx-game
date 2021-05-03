package ui

import (
	//"log"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/spriteloader"
	//"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/world"
)

type TilePaleteSelector struct {
	// store tiles for (1) displaying selector and (2) placing tiles
	Tiles []world.Tile
	Transform mgl32.Mat4
	Width int
	SelectedTileIndex int
	visible bool
}

func MakeTilePaleteSelector(tiles []world.Tile) TilePaleteSelector {
	width := math.Ceil(math.Sqrt(float64(len(tiles))))
	scale := float32(3/width)
	return TilePaleteSelector {
		Tiles: tiles,
		Transform: mgl32.Scale3D(scale,scale,scale),
		Width: int(width),
		SelectedTileIndex: -1,
	}
}

func (selector *TilePaleteSelector) Draw() {
	if !selector.visible {
		return
	}
	numTiles := float64(len(selector.Tiles))
	if numTiles>0 {
		for idx,tile := range selector.Tiles {
			yLocal := float32(idx/int(selector.Width))
			xLocal := float32(idx%int(selector.Width))
			localTransform := mgl32.Mat4.Mul4(
				selector.Transform,
				mgl32.Translate3D(xLocal,yLocal,0),
			)
			localPos := localTransform.Col(3)
			scaleX,scaleY,_ := mgl32.Extract3DScale(localTransform)
			// TODO add a custom texture for drawing air
			if tile.TileType!=world.TileTypeNone {
				spriteloader.DrawSpriteQuad(
					localPos.X(),localPos.Y(),
					scaleX,scaleY,
					int(tile.SpriteID),
				)
			}
		}
	}
}

func convertScreenCoordsToWorld(x,y float32, projection mgl32.Mat4) mgl32.Vec4 {
	homogenousClipCoords := mgl32.Vec4 { x,y,-1.0,1.0}
	cameraCoords := projection.Inv().Mul4x1(homogenousClipCoords)
	rayEye := mgl32.Vec4 { cameraCoords.X(), cameraCoords.Y(), -1.0, 0 }
	worldCoords := rayEye.Normalize().Mul(spriteloader.SpriteRenderDistance)
	worldCoords[3]=1
	return worldCoords
}

func (selector *TilePaleteSelector) TrySelectTile(x,y float32, projection mgl32.Mat4) bool {
	// can't select palete tile when palete is invisible
	if !selector.visible {
		return false
	}
	worldCoords := convertScreenCoordsToWorld(x,y,projection)
	paleteCoords := selector.Transform.Inv().Mul4x1(worldCoords).Vec2()
	tileX := int(math.Floor(float64(paleteCoords.X()+0.5)))
	tileY := int(math.Floor(float64(paleteCoords.Y()+0.5)))

	if tileX>=0 && tileX<selector.Width && tileY>=0 && tileY<selector.Width {
		selector.SelectedTileIndex = tileY*selector.Width + tileX
		return true
	}
	return false
}

func (selector *TilePaleteSelector) GetSelectedTile() world.Tile {
	if selector.SelectedTileIndex>=0 {
		return selector.Tiles[selector.SelectedTileIndex]
	} else {
		return world.Tile {}
	}
}


func (selector *TilePaleteSelector) Toggle() {
	selector.visible = !selector.visible
}
