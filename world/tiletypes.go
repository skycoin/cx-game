package world

import (
	"log"
	"io/ioutil"

	"github.com/go-yaml/yaml"

	"github.com/skycoin/cx-game/config"
	"github.com/skycoin/cx-game/engine/spriteloader/blobsprites"
	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/render"
)

type TileConfig struct {
	Blob         string   `yaml:"blob"`
	Collision    string   `yaml:"collision"`
	Sprites      []string `yaml:"sprites"`
	Sprite       string   `yaml:"sprite"`
	Layer        string   `yaml:"layer"`
	Invulnerable bool     `yaml:"invulnerable"`
	Category     string   `yaml:"category"`
	Width        int32    `yaml:"width"`
	Height       int32    `yaml:"height"`
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

func (config *TileConfig) Placer(fname string, id TileTypeID) Placer {

	if config.Blob == "" {
		// TODO handle multiple sprites
		spriteID := render.GetSpriteIDByName(config.Sprite)
		return DirectPlacer {
			SpriteID: spriteID,
			Category: tileCategoryFromString(config.Category),
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

func (config *TileConfig) TileType(name string, id TileTypeID) TileType {
	return TileType{
		Name:         name,
		Layer:        LayerIDFromName(config.Layer),
		Placer:       config.Placer(name, id),
		Invulnerable: config.Invulnerable,
		Width:        config.Width, Height: config.Height,
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
	})
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
			RegisterTileType(name, config.TileType(name, NextTileTypeID()))
		}
	}
}
