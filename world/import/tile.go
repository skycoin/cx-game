package worldimport

import (
	"fmt"
	"image"
	"path"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/lafriks/go-tiled"

	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/render"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func defaltToolForLayer(layerID world.LayerID) types.ToolType {
	if layerID == world.BgLayer {
		return constants.BG_TOOL
	}
	return constants.FURNITURE_TOOL
}

func registerTilesetTile(
	imgPath string, id uint32,
	transform mgl32.Mat3, model mgl32.Mat4, layerID world.LayerID,
) world.TileTypeID {
	tex := spriteloader.LoadTextureFromFileToGPUCached(imgPath)
	name := fmt.Sprintf("%v:%v", imgPath, id)
	sprite := render.Sprite{
		Name:      name,
		Transform: transform,
		Model:     model,
		Texture:   render.Texture{Target: gl.TEXTURE_2D, Texture: tex.Gl},
	}
	spriteID := render.RegisterSprite(sprite)

	mw := int32(math32.Round(model.At(0, 0)))
	mh := int32(math32.Round(model.At(1, 1)))

	tile := world.NewNormalTile()
	tile.Name = name
	tile.TileTypeID = world.NextTileTypeID()

	tileType := world.TileType{
		Name:   name,
		Layer:  layerID,
		Placer: world.DirectPlacer{SpriteID: spriteID, Tile: tile},
		Width:  mw, Height: mh,
	}

	tileTypeID :=
		world.RegisterTileType(name, tileType, defaltToolForLayer(layerID))

	return tileTypeID
}

func findTilesetTileForLayerTile(
	tilesetTiles []*tiled.TilesetTile, layerTile *tiled.LayerTile,
) (*tiled.TilesetTile, bool) {
	for _, tilesetTile := range layerTile.Tileset.Tiles {
		if tilesetTile.ID == layerTile.ID {
			return tilesetTile, true
		}
	}
	return nil, false
}

// compute a 3x3 matrix which maps every point in the
// "parent" rectangle to a corresponding point in the "here" rectangle
func rectTransform(here image.Rectangle, parentDims image.Point) mgl32.Mat3 {
	scaleX := float32(here.Dx()) / float32(parentDims.X)
	scaleY := float32(here.Dy()) / float32(parentDims.Y)
	scale := mgl32.Scale2D(scaleX, scaleY)

	translate := mgl32.Translate2D(
		float32(here.Min.X)/float32(parentDims.X),
		float32(here.Min.Y)/float32(parentDims.Y))

	return translate.Mul3(scale)
}

func modelFromSize(dx int, dy int) mgl32.Mat4 {
	return mgl32.Scale2D(float32(dx)/16, float32(dy)/16).Mat4()
}

type TilesetIDKey struct {
	tileset *tiled.Tileset
	id      uint32
	scaleX  int
	scaleY  int
}

var tilesetAndIDToCXTile = map[TilesetIDKey]world.TileTypeID{}

func getTileTypeID(
	layerTile *tiled.LayerTile, tmxPath string, layerID world.LayerID,
) world.TileTypeID {
	tileset := layerTile.Tileset
	// nil entry => empty layer tile
	if tileset == nil {
		return world.TileTypeIDAir
	}

	// search for tile in existing tiles
	tilesetTile, foundTilesetTile :=
		findTilesetTileForLayerTile(tileset.Tiles, layerTile)

	if foundTilesetTile {
		cxtile := tilesetTile.Properties.GetString("cxtile")
		tileTypeID, foundTileTypeID := world.IDFor(cxtile)
		if foundTileTypeID {
			return tileTypeID
		}
	}

	flipX,flipY := scaleFromFlipFlags(layerTile)
	flipTransform := mgl32.Scale2D( float32(flipX), float32(flipY) )
	key := TilesetIDKey{tileset, layerTile.ID, flipX, flipY}
	cachedTileTypeID, hitCache := tilesetAndIDToCXTile[key]
	if hitCache {
		return cachedTileTypeID
	}

	// did not find - register new tile type
	if foundTilesetTile && tilesetTile.Image != nil {
		imgPath := path.Join(tmxPath, "..", tilesetTile.Image.Source)
		w := tilesetTile.Image.Width
		h := tilesetTile.Image.Height
		model := modelFromSize(w, h)
		tileTypeID := registerTilesetTile(
			imgPath, layerTile.ID, flipTransform, model, layerID)
		tilesetAndIDToCXTile[key] = tileTypeID
		return tileTypeID
	}
	// register new tile from tileset
	imgPath := path.Join(tmxPath, "..", tileset.Image.Source)
	rect := tileset.GetTileRect(layerTile.ID)
	transform := rectTransform(rect,
		image.Point{tileset.Image.Width, tileset.Image.Height},
	)
	transform = transform.Mul3(flipTransform)
	tileTypeID := registerTilesetTile(
		imgPath, layerTile.ID, transform, mgl32.Ident4(), layerID)
	tilesetAndIDToCXTile[key] = tileTypeID
	return tileTypeID
}
