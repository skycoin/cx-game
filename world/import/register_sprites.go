package worldimport

import (
	"fmt"
	"log"

	"github.com/lafriks/go-tiled"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
)

// a sprite registered from a tiled import
type TiledSprite struct {
	Image    TilesetTileImage
	Metadata TiledMetadata
}

func (ts TiledSprite) Register(name string) RegisteredTiledSprite {
	return RegisteredTiledSprite{
		SpriteID: ts.Image.RegisterSprite(name),
		Metadata: ts.Metadata, // just copy metadata over
		Width:    ts.Image.Width,
		Height:   ts.Image.Height,
	}
}

type RegisteredTiledSprite struct {
	SpriteID render.SpriteID
	Width    int32
	Height   int32
	Metadata TiledMetadata
}

// properties on a Tiled tileset tile that are relevant to cx-game
type TiledMetadata struct {
	Powered OptionalBool
	Name    string
	LayerID world.LayerID
}

func NewTiledMetadata(name string) TiledMetadata {
	return TiledMetadata{Name: name, LayerID: world.TopLayer}
}

type OptionalBool struct {
	Set   bool
	Value bool
}

type LayerTiledSpritePair struct {
	LayerID world.LayerID
	Sprite  TiledSprite
}

type TiledSprites map[string][]TiledSprite
type RegisteredTiledSprites map[string][]RegisteredTiledSprite

func findTiledSpritesInMapTilesets(
	tiledMap *tiled.Map, mapDir string,
) TiledSprites {
	tiledSprites := TiledSprites{}

	for _, tileset := range tiledMap.Tilesets {
		registeredTileIDs := map[uint32]bool{}
		for _, tilesetTile := range tileset.Tiles {
			name := nameForTilesetTile(tileset.Name, tilesetTile.ID)
			metadata := NewTiledMetadata(name)
			metadata.ParseFrom(tilesetTile.Properties)
			image := imageForTilesetTile(
				tileset, tilesetTile.ID, tilesetTile, mapDir)
			tiledSprite := TiledSprite{Image: image, Metadata: metadata}
			tiledSprites[metadata.Name] =
				append(tiledSprites[metadata.Name], tiledSprite)
			registeredTileIDs[tilesetTile.ID] = true
		}
		if tileset.Image != nil {
			tileCount := uint32(
				tileset.Image.Width / tileset.TileWidth *
					tileset.Image.Height / tileset.TileHeight)

			for id := uint32(0); id < tileCount; id++ {
				name := nameForTilesetTile(tileset.Name, id)
				metadata := NewTiledMetadata(name)
				isRegistered := registeredTileIDs[id]
				if !isRegistered {
					image :=
						imageForTilesetTile(tileset, uint32(id), nil, mapDir)
					tiledSprite :=
						TiledSprite{Image: image, Metadata: metadata}
					tiledSprites[metadata.Name] =
						append(tiledSprites[metadata.Name], tiledSprite)
				}
			}
		}
	}

	return tiledSprites
}

func registerTiledSprites(tiledSprites TiledSprites) RegisteredTiledSprites {
	registeredTiledSprites := RegisteredTiledSprites{}
	for name, tileSprites := range tiledSprites {
		registeredTiledSprites[name] = []RegisteredTiledSprite{}
		for _, tileSprite := range tileSprites {
			_, ok := world.IDFor(name)
			if !ok {
				registeredTiledSprites[name] = append(
					registeredTiledSprites[name], tileSprite.Register(name))
			}
		}
	}
	return registeredTiledSprites
}

func hasProperty(properties tiled.Properties, name string) bool {
	for _, property := range properties {
		if property.Name == name {
			return true
		}
	}
	return false
}

func (metadata *TiledMetadata) ParseFrom(properties tiled.Properties) {
	if hasProperty(properties, "powered") {
		metadata.Powered.Set = true
		metadata.Powered.Value = properties.GetBool("powered")
	}
	if hasProperty(properties, "cxtile") {
		metadata.Name = properties.GetString("cxtile")
	}
}

func parseMetadataFromLayerTile(layerTile *tiled.LayerTile) TiledMetadata {
	name := fmt.Sprintf("%v:%v", layerTile.Tileset.Name, layerTile.ID)
	metadata := NewTiledMetadata(name)
	tilesetTile, ok := findTilesetTileForLayerTile(layerTile)
	if ok {
		metadata.ParseFrom(tilesetTile.Properties)
	}
	return metadata
}
