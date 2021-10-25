package world

import (
	"log"
	"strings"

	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/spriteloader"
)

type LayerID int

const (
	BgLayer LayerID = iota
	WindowLayer
	PipeLayer
	MidLayer
	TopLayer

	NumLayers // DO NOT SET MANUALLY
)

func (id LayerID) Valid() bool {
	return id >= 0 && id < NumLayers
}

type Layer struct {
	Tiles       []Tile
	Spritesheet spriteloader.Spritesheet
}

type Layers [NumLayers]Layer

func NewLayer(numTiles int32) Layer {
	return Layer{
		Tiles: make([]Tile, numTiles),
	}
}

func NewLayers(numTiles int32) Layers {
	layers := Layers{}
	for i := LayerID(0); i < NumLayers; i++ {
		layers[i] = NewLayer(numTiles)
	}
	return layers
}

var layerIDsByName = map[string]LayerID{
	"foreground": TopLayer,
	"objects":    MidLayer,
	"windows":    WindowLayer,
	"walls":      BgLayer,
	"pipesim":    PipeLayer,
}

func LayerIDForName(name string) (LayerID, bool) {
	name = strings.ToLower(name)
	layerID, ok := layerIDsByName[name]
	return layerID, ok
}

func (layerID LayerID) Z() float32 {
	if layerID == TopLayer {
		return constants.FRONTLAYER_Z
	} else if layerID == MidLayer {
		return constants.MIDLAYER_Z
	} else if layerID == BgLayer {
		return constants.BGLAYER_Z
	} else if layerID == PipeLayer {
		return constants.PIPELAYER_Z
	} else if layerID == WindowLayer {
		return constants.WINDOWLAYER_Z
	} else {
		log.Fatalf("error: Unknown layer!")
	}
	return -1
}
