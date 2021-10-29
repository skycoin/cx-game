package worldimport

import (
	"fmt"
	"image"
	"log"
	"path"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lafriks/go-tiled"
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
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

	key := TilesetIDKey{tileset, layerTile.ID}
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
		//transform := mgl32.Scale2D(16/float32(w), 16/float32(h))
		transform := mgl32.Ident3()
		tileTypeID := registerTilesetTile(
			imgPath, layerTile.ID, transform, model, layerID)
		tilesetAndIDToCXTile[key] = tileTypeID
		return tileTypeID
	}
	// register new tile from tileset
	imgPath := path.Join(tmxPath, "..", tileset.Image.Source)
	rect := tileset.GetTileRect(layerTile.ID)
	transform := rectTransform(rect,
		image.Point{tileset.Image.Width, tileset.Image.Height},
	)
	tileTypeID := registerTilesetTile(
		imgPath, layerTile.ID, transform, mgl32.Ident4(), layerID)
	tilesetAndIDToCXTile[key] = tileTypeID
	return tileTypeID
}

func importTile(
	planet *world.Planet,
	tileIndex int, layerTile *tiled.LayerTile, tmxPath string,
	layerID world.LayerID,
) {
	tileTypeID := getTileTypeID(layerTile, tmxPath, layerID)
	if tileTypeID != world.TileTypeIDAir {

		// correct mismatch between Tiled Y axis (downwards)
		// and our Y axis  (upwards)
		y := int(planet.Height) - tileIndex/int(planet.Width)
		x := tileIndex % int(planet.Width)
		planet.PlaceTileType(tileTypeID, x, y)
	}
}

func importLayer(
	planet *world.Planet, tiledLayer *tiled.Layer, tmxPath string,
	layerID world.LayerID,
) {
	for idx, layerTile := range tiledLayer.Tiles {
		importTile(planet, idx, layerTile, tmxPath, layerID)
	}
}

func ImportWorld(tmxPath string) world.World {
	start := time.Now()
	tiledMap, err := tiled.LoadFromFile(tmxPath)
	if err != nil {
		log.Fatalf("import world: %v", err)
	}
	elapsedTiledLoad := time.Since(start)
	log.Printf("load %s took %s", tmxPath, elapsedTiledLoad)
	planet := world.NewPlanet(int32(tiledMap.Width), int32(tiledMap.Height))
	for _, tiledLayer := range tiledMap.Layers {
		layerID, foundLayerID := world.LayerIDForName(tiledLayer.Name)
		if foundLayerID {
			importLayer(planet, tiledLayer, tmxPath, layerID)
		}
	}
	return world.World{
		Planet: *planet,
		Entities: world.Entities{
			Agents: *agents.NewAgentList(),
		},
		Stats:     world.NewWorldStats(),
		TimeState: world.NewTimeState(),
	}
}
