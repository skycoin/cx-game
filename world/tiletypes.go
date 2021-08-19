package world

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-yaml/yaml"

	"github.com/skycoin/cx-game/config"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/spriteloader/blobsprites"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/render/blob"
)

type TileConfig struct {
	Blob         string   `yaml:"blob"`
	Collision    string   `yaml:"collision"`
	Sprites      []string `yaml:"sprites"`
	Sprite       string   `yaml:"sprite"`
	Layer        string   `yaml:"layer"`
	Invulnerable bool     `yaml:"invulnerable"`
	Category     string   `yaml:"category"`
	/*
		Width        int32    `yaml:"width"`
		Height       int32    `yaml:"height"`
	*/
}

func loadIDsFromSpritenames(names []string, n int) []blobsprites.BlobSpritesID {
	ids := make([]blobsprites.BlobSpritesID, len(names))
	for idx, name := range names {
		ids[idx] = blobsprites.LoadIDFromSpritename(name, n)

	}
	return ids
}

func tileCategoryFromString(str string) TileCategory {
	if str == "" {
		return TileCategoryNormal
	}
	if str == "liquid" {
		return TileCategoryLiquid
	}

	log.Fatalf("unrecognized tile category [%v]", str)
	return TileCategoryNone
}

func tileCollisionTypeFromString(str string) TileCollisionType {
	if str == "" {
		return TileCollisionTypeSolid
	}
	if str == "platform" {
		return TileCollisionTypePlatform
	}

	log.Fatalf("unrecognized tile collision type [%v]", str)
	return TileCollisionTypeSolid
}

func (config *TileConfig) Placer(fname string, id TileTypeID) Placer {

	if config.Blob == "" {
		// TODO handle multiple sprites
		spriteID := render.GetSpriteIDByName(config.Sprite)
		return DirectPlacer{
			SpriteID:          spriteID,
			Category:          tileCategoryFromString(config.Category),
			TileCollisionType: tileCollisionTypeFromString(config.Collision),
		}
	}
	if config.Blob == "full" {
		ids := loadIDsFromSpritenames(config.Sprites,
			blob.BlobSheetWidth*blob.BlobSheetHeight)

		return AutoPlacer{
			blobSpritesIDs: ids,
			TileTypeID:     id, TilingType: blob.FullBlobTiling,
		}
	}
	if config.Blob == "simple" {
		ids := loadIDsFromSpritenames(
			config.Sprites,
			blob.SimpleBlobSheetWidth*blob.SimpleBlobSheetHeight)

		return AutoPlacer{
			blobSpritesIDs: ids,
			TileTypeID:     id, TilingType: blob.SimpleBlobTiling,
		}
	}

	log.Fatalf("unrecognized blob type: %s", config.Blob)
	return DirectPlacer{}
}

var layerNamesToIDs = map[string]LayerID{
	"top": TopLayer, "mid": MidLayer, "bg": BgLayer,
}

func LayerIDFromName(name string) LayerID {
	id, ok := layerNamesToIDs[name]
	if !ok {
		log.Fatalf("unknown layer name [%v]", name)
	}
	return id
}

func (config *TileConfig) Dims() (int32, int32) {
	if config.Sprite == "" {
		return 1, 1
	}
	spriteID := render.GetSpriteIDByName(config.Sprite)
	sprite := render.GetSpriteByID(spriteID)
	model := sprite.Model

	width := int32(model.At(0, 0))
	height := int32(model.At(1, 1))
	return width, height
}

func (config *TileConfig) TileType(name string, id TileTypeID) TileType {
	width, height := config.Dims()
	return TileType{
		Name:         name,
		Layer:        LayerIDFromName(config.Layer),
		Placer:       config.Placer(name, id),
		Invulnerable: config.Invulnerable,
		Width:        width, Height: height,
	}
}

type TileConfigs map[string]TileConfig

func RegisterTileTypes() {
	RegisterEmptyTileType()
	RegisterConfigTileTypes()
}

func RegisterEmptyTileType() {
	RegisterTileType("air", TileType{
		Name:   "Air",
		Placer: DirectPlacer{},
	}, constants.NULL_TOOL)
}

func RegisterConfigTileTypes() {
	paths := config.FindPathsWithExt("./assets/tile/", ".yaml")

	for _, path := range paths {
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("cannot read tile config at %s", path)
		}
		var configs TileConfigs
		err = yaml.Unmarshal(buf, &configs)
		if err != nil {
			log.Fatalf("parse tile config %s: %v", path, err)
		}
		for name, config := range configs {
			//todo this is quick hack
			if strings.Contains(path, "tiles.yaml") {
				RegisterTileType(name, config.TileType(name, NextTileTypeID()), constants.TILE_TOOL)
			} else {
				RegisterTileType(name, config.TileType(name, NextTileTypeID()), constants.FURNITURE_TOOL)
			}
		}
	}
}
