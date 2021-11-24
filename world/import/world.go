package worldimport

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/lafriks/go-tiled"
	"github.com/skycoin/cx-game/world"
)

func importTile(
	planet *world.Planet,
	tileIndex int, layerTile *tiled.LayerTile, tmxPath string,
	layerID world.LayerID, tileTypeIDs map[string]world.TileTypeID,
) {
	if layerTile.Nil {
		return
	}
	// correct mismatch between Tiled Y axis (downwards)
	// and our Y axis  (upwards)
	y := int(planet.Height) - tileIndex/int(planet.Width)
	x := tileIndex % int(planet.Width)

	name := nameForLayerTile(layerTile)
	tileTypeID, ok := world.IDFor(name)
	if !ok {
		tileTypeID, ok = tileTypeIDs[name]
		if !ok {
			log.Fatalf("cannot found tile type ID for %v at %d,%d", name, x,y)
		}
	}
	if tileTypeID != world.TileTypeIDAir {
		opts := world.NewTileCreationOptions()
		flipX, flipY := scaleFromFlipFlags(layerTile)
		opts.FlipTransform = mgl32.Scale3D(float32(flipX), float32(flipY), 1)
		planet.PlaceTileType(tileTypeID, x, y, opts)
	}
}

func importLayer(
	planet *world.Planet, tiledLayer *tiled.Layer, tmxPath string,
	layerID world.LayerID, tileTypeIDs map[string]world.TileTypeID,
) {
	for idx, layerTile := range tiledLayer.Tiles {
		importTile(planet, idx, layerTile, tmxPath, layerID, tileTypeIDs)
	}
}

func filterTiledSpritesInMapLayers(
	allTiledSprites TiledSprites, tiledMap *tiled.Map,
) TiledSprites {
	mapTiledSprites := TiledSprites{}
	for _, layer := range tiledMap.Layers {
		for _, layerTile := range layer.Tiles {
			if !layerTile.Nil {
				name := nameForLayerTile(layerTile)
				layerID, _ := world.LayerIDForName(layer.Name)
				mapTiledSprites[name] = allTiledSprites[name]
				for idx := range mapTiledSprites[name] {
					mapTiledSprites[name][idx].Metadata.LayerID = layerID
				}
			}
		}
	}

	return mapTiledSprites
}
