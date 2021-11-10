package worldimport

import (
	"image"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/lafriks/go-tiled"

	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
)

func defaltToolForLayer(layerID world.LayerID) types.ToolType {
	if layerID == world.BgLayer {
		return constants.BG_TOOL
	}
	return constants.FURNITURE_TOOL
}

func findTilesetTileForLayerTile(
	layerTile *tiled.LayerTile,
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

type TilesetIDKey struct {
	tileset *tiled.Tileset
	id      uint32
}

var tilesetAndIDToCXTile = map[TilesetIDKey]world.TileTypeID{}

func directPlacerForTileSprites(tileSprites []RegisteredTiledSprite) world.Placer {
	tile := world.NewNormalTile()
	tile.Name = tileSprites[0].Metadata.Name
	tile.TileTypeID = world.NextTileTypeID()

	return world.DirectPlacer{
		SpriteID: tileSprites[0].SpriteID, Tile: tile,
	}
}

func powerPlacerForTileSprites(tileSprites []RegisteredTiledSprite) world.Placer {
	tile := world.NewNormalTile()
	tile.Name = tileSprites[0].Metadata.Name
	tile.TileTypeID = world.NextTileTypeID()

	placer := world.LightPlacer{Tile: tile}
	for _, tileSprite := range tileSprites {
		if tileSprite.Metadata.Powered.Value {
			placer.OnSpriteID = tileSprite.SpriteID
		} else {
			placer.OffSpriteID = tileSprite.SpriteID
		}
	}
	return placer
}

func placerForTileSprites(tileSprites []RegisteredTiledSprite) world.Placer {
	// use powered placer if "powered" field is set on any relevant sprites
	for _, tileSprite := range tileSprites {
		if tileSprite.Metadata.Powered.Set {
			return powerPlacerForTileSprites(tileSprites)
		}
	}
	// otherwise, use direct placer (1:1 sprite:tiletype ratio)
	return directPlacerForTileSprites(tileSprites)
}

func registerTileTypeForTileSprites(
	tileSprites []RegisteredTiledSprite,
) world.TileTypeID {
	layerID := tileSprites[0].Metadata.LayerID
	name := tileSprites[0].Metadata.Name
	tileType := world.TileType{
		Name:  name,
		Width: tileSprites[0].Width, Height: tileSprites[0].Height,
		Placer: placerForTileSprites(tileSprites),
		Layer:  layerID,
	}

	tileTypeID :=
		world.RegisterTileType(name, tileType, defaltToolForLayer(layerID))
	return tileTypeID
}

func registerTileTypesForTiledSprites(
	tiledSprites RegisteredTiledSprites,
) map[string]world.TileTypeID {
	tileTypeIDs := map[string]world.TileTypeID{}

	for name, tileSprites := range tiledSprites {
		if len(tileSprites)>0 {
			tileTypeIDs[name] = registerTileTypeForTileSprites(tileSprites)
		}
	}

	return tileTypeIDs
}

func getTileTypeID(
	layerTile *tiled.LayerTile, tmxPath string, layerID world.LayerID,
	tiledSprites TiledSprites,
) world.TileTypeID {
	tileset := layerTile.Tileset
	// nil entry => empty layer tile
	if tileset == nil {
		return world.TileTypeIDAir
	}

	// search for tile in existing tiles
	tilesetTile, foundTilesetTile := findTilesetTileForLayerTile(layerTile)

	if foundTilesetTile {
		cxtile := tilesetTile.Properties.GetString("cxtile")
		tileTypeID, foundTileTypeID := world.IDFor(cxtile)
		if foundTileTypeID {
			return tileTypeID
		}
	}

	flipX, flipY := scaleFromFlipFlags(layerTile)
	flipTransform := mgl32.Scale2D(float32(flipX), float32(flipY))
	key := TilesetIDKey{tileset, layerTile.ID}
	cachedTileTypeID, hitCache := tilesetAndIDToCXTile[key]
	if hitCache {
		return cachedTileTypeID
	}

	// did not find - register new tile type
	return registerTilesetTile(layerTile, TileRegistrationOptions{
		TmxPath: tmxPath, LayerID: layerID, Tileset: tileset,
		LayerTile: layerTile, TilesetTile: tilesetTile,
		FlipTransform: flipTransform,
		TiledSprites:  tiledSprites,
	})
}
