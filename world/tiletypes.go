package world

import (
	"os"
	"log"

	"github.com/go-yaml/yaml"
	"github.com/skycoin/cx-game/render/blob"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/spriteloader/blobsprites"
)

type TileConfig struct {
	Blob string `yaml:"blob"`
	Collision string `yaml:"collision"`
	Sprites []string `yaml:"sprites"`
	Sprite string `yaml:"sprite"`
	Layer string `yaml:"layer"`
	Invulnerable bool `yaml:"invulnerable"`
}

func loadIDsFromSpritenames(names []string) []blobsprites.BlobSpritesID {
	ids := make([]blobsprites.BlobSpritesID, len(names))
	for idx,name := range names {
		// TODO support simple tiling also
		ids[idx] = blobsprites.LoadIDFromSpritename(
			name, blob.BlobSheetWidth * blob.BlobSheetHeight )

	}
	return ids
}

func (config *TileConfig) Placer(fname string, id TileTypeID) Placer {

	if config.Blob == "" {
		// TODO handle multiple sprites
		spriteID := spriteloader.GetSpriteIdByName(config.Sprite)
		return DirectPlacer { SpriteID: spriteID }
	}
	if config.Blob == "full" {
		ids := loadIDsFromSpritenames(config.Sprites)
		return AutoPlacer {
			// TODO
			blobSpritesIDs: ids,
			TileTypeID: id, TilingType: blob.FullBlobTiling,
		}
	}

	log.Fatalf("unrecognized blob type: %s",config.Blob)
	return DirectPlacer{}
}

var layerNamesToIDs = map[string] LayerID {
	"top": TopLayer, "mid": MidLayer, "bg": BgLayer,
}

func LayerIDFromName(name string) LayerID{
	id,ok := layerNamesToIDs[name]
	if !ok { log.Fatalf("unknown layer name [%v]", name) }
	return id
}

func (config *TileConfig) TileType(name string, id TileTypeID) TileType {
	return TileType {
		Name: name,
		Layer: LayerIDFromName(config.Layer),
		Placer: config.Placer(name,id),
		Invulnerable: config.Invulnerable,
	}
}

type TileConfigs map[string]TileConfig


func RegisterTileTypes() {
	RegisterEmptyTileType()
	RegisterConfigTileTypes()
}

func RegisterEmptyTileType() {
	RegisterTileType("air", TileType {
		Name: "Air",
		Placer: DirectPlacer{},
	})
}

const tilesConfigPath = "./assets/tile/tiles.yaml"

func RegisterConfigTileTypes() {
	buf, err := os.ReadFile(tilesConfigPath)
	if err != nil {
		log.Fatalf("cannot read tile config at %s",tilesConfigPath)
	}
	var configs TileConfigs
	err = yaml.Unmarshal(buf,&configs)
	if err != nil {
		log.Fatalf("parse tile configs: %v",err)
	}
	for name,config := range configs {
		RegisterTileType(name, config.TileType(name,NextTileTypeID()))
	}

}
