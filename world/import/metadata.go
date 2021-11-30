package worldimport

import (
	"fmt"

	"github.com/lafriks/go-tiled"

	"github.com/skycoin/cx-game/world"
)

// properties on a Tiled tileset tile that are relevant to cx-game
type TiledMetadata struct {
	Powered OptionalBool
	Name    string
	LayerID world.LayerID
	NeedsGround bool
	Wattage int
}

func NewTiledMetadata(name string) TiledMetadata {
	return TiledMetadata{Name: name, LayerID: world.TopLayer}
}

type OptionalBool struct {
	Set   bool
	Value bool
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
	metadata.NeedsGround = properties.GetBool("needsground")
	metadata.Wattage = properties.GetInt("wattage")
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
