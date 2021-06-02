package ui

import (
	"log"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/render"
)

type layerIndexPair struct {layer, index int}

type TilePaletteSelector struct {
	// store tiles for (1) displaying selector and (2) placing tiles
	Layers            []PaleteLayer
	Transform         mgl32.Mat4
	Width             int
	Height            int
	SelectedTileIndex int
	LayerIndex        int
	bluePixelSpriteId int
	redPixelSpriteId  int
	multiTiles        map[layerIndexPair]world.MultiTile
}

type PaleteLayer []world.Tile

func NewPaleteLayers(slots int) []PaleteLayer {
	layers := make([]PaleteLayer,numLayers)
	for layer:=0; layer<numLayers; layer++ {
		layers[layer] = make([]world.Tile,slots)
		for slot:=0; slot<slots; slot++ {
			layers[layer][slot] = world.NewEmptyTile()
		}
	}
	return layers
}

const numLayers = 3
var layerNames = []string { "Background","Mid","Top" }

func (selector *TilePaletteSelector) CycleLayer() {
	selector.LayerIndex++
	if selector.LayerIndex == numLayers {
		selector.LayerIndex = -1
	}
}

func (selector *TilePaletteSelector) Visible() bool {
	return selector.LayerIndex >= 0
}

func (selector *TilePaletteSelector) SelectedTile() *world.Tile {
	return &selector.Tiles()[selector.SelectedTileIndex]
}

func (selector *TilePaletteSelector) IsMultiTileSelected() bool {
	if !selector.Visible() { return false }
	return selector.SelectedTile().TileType == world.TileTypeMulti
}

const selectorSize = 5

func MakeTilePaleteSelector(width,height int) TilePaletteSelector {

	spriteloader.LoadSingleSprite("./assets/blue_pixel.png", "blue-pixel")
	spriteloader.LoadSingleSprite("./assets/red_pixel.png", "red-pixel")

	return TilePaletteSelector{
		Layers:            NewPaleteLayers(width*height),
		Width:             width,
		Height:            height,
		SelectedTileIndex: -1,
		LayerIndex:        -1,
		bluePixelSpriteId: spriteloader.GetSpriteIdByName("blue-pixel"),
		redPixelSpriteId:  spriteloader.GetSpriteIdByName("red-pixel"),
		multiTiles:        make(map[layerIndexPair]world.MultiTile),
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

	if !selector.Visible() {
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

	tiles := selector.Layers[selector.LayerIndex]
	if len(tiles) > 0 {
		for idx, tile := range tiles {
			tileWorldTransform := selector.worldTransformForTileAtIndex(idx).
				Mul4(mgl32.Scale3D(1-spacing, 1-spacing, 1))

			// TODO add a custom texture for drawing air
			if tile.TileType != world.TileTypeNone {
				spriteloader.DrawSpriteQuadMatrix(tileWorldTransform, int(tile.SpriteID))
			}
		}
	}

	textCtx := selectorCtx.
		PushLocal(mgl32.Translate3D(0,1,0).Mul4(cxmath.Scale(0.4)))
	DrawString(
		layerNames[selector.LayerIndex], mgl32.Vec4{1,1,1,1},
		AlignLeft, textCtx,
	)

}

func (selector *TilePaletteSelector) TrySelectTile(x, y float32) bool {
	// can't select palete tile when palete is invisible
	if !selector.Visible() {
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
		childTile := selector.Tiles()[tileY*selector.Width + tileX]
		parentX := tileX - int(childTile.OffsetX)
		parentY := tileY - int(childTile.OffsetY)
		parentIdx := parentY*selector.Width + parentX
		selector.SelectedTileIndex = parentIdx
		return true
	}
	return false
}

func (selector *TilePaletteSelector) Tiles() []world.Tile {
	return selector.Layers[selector.LayerIndex]
}

func (selector *TilePaletteSelector) GetSelectedTile() world.Tile {
	if selector.SelectedTileIndex >= 0 &&
		selector.SelectedTileIndex < len(selector.Tiles()) {

		return selector.Tiles()[selector.SelectedTileIndex]
	} else {
		return world.Tile{}
	}
}

func (selector *TilePaletteSelector) AddTile(
		tile world.Tile,
		x,y int,
		layerIndex world.Layer,
) {
	idx := y*selector.Width + x
	layer := selector.Layers[layerIndex]
	layer[idx] = tile
}

func (selector *TilePaletteSelector) AddMultiTile(
		multiTile world.MultiTile,
		left,top int,
		layerIndex world.Layer,
) {
	for _,tile := range multiTile.Tiles() {
		selector.AddTile(
			tile,
			left+int(tile.OffsetX),
			top+int(tile.OffsetY),
			layerIndex,
		)
	}
	tileIdx := top*selector.Width + left
	pair := layerIndexPair {int(layerIndex), tileIdx}
	selector.multiTiles[pair] = multiTile
}

func (selector *TilePaletteSelector) GetSelectedMultiTile() world.MultiTile {
	pair := layerIndexPair { selector.LayerIndex, selector.SelectedTileIndex }
	mt,ok := selector.multiTiles[pair]
	if !ok {
		log.Fatal("Tile palette selector cannot find selected multi-tile")
	}
	return mt
}
