package worldimport

import (
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
