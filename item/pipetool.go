package item

// TODO where to store state about previous mouse position

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/world"
	"github.com/skycoin/cx-game/world/pipesim"
)

func pipetoolMouseDownRight(info ItemUseInfo) bool {
	// always consume
	return true
}

func pipetoolMouseDown(info ItemUseInfo) {
	didSelect := info.Inventory.PlacementGrid.TrySelect(info.CamCoords())
	if didSelect {
		return
	}
	didPlace := info.Inventory.PlacementGrid.TryPlaceNoConnect(info)
	_ = didPlace
}

func pipetoolMouseDrag(
		info ItemUseInfo, prevMousePos mgl32.Vec2, b glfw.MouseButton,
) {
	currMousePos := info.WorldCoords()

	prevMouseTile := cxmath.TileAt(prevMousePos)
	currMouseTile := cxmath.TileAt(currMousePos)

	// place new tile
	didPlace := false
	if b == glfw.MouseButtonLeft {
		didPlace = info.Inventory.PlacementGrid.TryPlaceNoConnect(info)
	}

	planet := &info.World.Planet
	pipeLayerTiles := planet.GetLayerTiles(world.PipeLayer)
	prevTileIdx :=
		planet.GetTileIndex(int(prevMouseTile.X), int(prevMouseTile.Y))
	currTileIdx :=
		planet.GetTileIndex(int(currMouseTile.X), int(currMouseTile.Y))

	prevTile := &pipeLayerTiles[prevTileIdx]
	currTile := &pipeLayerTiles[currTileIdx]

	if didPlace { currTile.Connections = pipesim.Connections{} }

	// abort if either involved tiles are empty
	if currTile.TileCategory == world.TileCategoryNone { return }
	if prevTile.TileCategory == world.TileCategoryNone { return }
	// don't try to connect tile to itself
	if currMouseTile == prevMouseTile { return }

	// connect prev tile to new tile

	disp := currMouseTile.Sub(prevMouseTile)
	prevNewConnections, currNewConnections := pipesim.FindNewConnections(disp)

	if b == glfw.MouseButtonLeft {
		prevTile.Connections = prevTile.Connections.OR(prevNewConnections)
		currTile.Connections = currTile.Connections.OR(currNewConnections)
	} else { // probably mouse button right
		prevTile.Connections =
			prevTile.Connections.AND(prevNewConnections.NOT())
		currTile.Connections =
			currTile.Connections.AND(currNewConnections.NOT())
	}

	currTile.TileTypeID.Get().UpdateTile(world.TileUpdateOptions{
		Tile:       currTile,
		Cycling:    true,
		Neighbours: planet.GetNeighbours(
			pipeLayerTiles, int(currMouseTile.X), int(currMouseTile.Y),
			currTile.TileTypeID,
		),
	})
	prevTile.TileTypeID.Get().UpdateTile(world.TileUpdateOptions{
		Tile:       prevTile,
		Cycling:    true,
		Neighbours: planet.GetNeighbours(
			pipeLayerTiles, int(prevMouseTile.X), int(prevMouseTile.Y),
			prevTile.TileTypeID,
		),
	})
}

func RegisterPipeToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-pipe-place-tool"))
	itemtype.Name = "Dev Pipe Place Tool"
	itemtype.Category = BuildTool
	itemtype.Use = pipetoolMouseDown
	itemtype.OnDrag = pipetoolMouseDrag
	itemtype.MouseDownRight = pipetoolMouseDownRight
	return AddItemType(itemtype)
}
