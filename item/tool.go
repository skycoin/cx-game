package item

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/cxmath"
)

func RegisterFurnitureToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-furniture-tool-2"))
	itemtype.Name = "Dev Furniture Tool"
	itemtype.Category = BuildTool
	itemtype.Use = UseBuildTool
	return AddItemType(itemtype)
}

func RegisterTileToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-tile-tool"))
	itemtype.Name = "Dev Tile Tool"
	itemtype.Category = BuildTool
	itemtype.Use = UseBuildTool
	return AddItemType(itemtype)
}

func RegisterBgToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-tile-tool"))
	itemtype.Name = "Dev Background Tile Tool"
	itemtype.Category = BuildTool
	itemtype.Use = UseBuildTool
	return AddItemType(itemtype)
}

func RegisterPipeToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-pipe-place-tool"))
	itemtype.Name = "Dev Pipe Place Tool"
	itemtype.Category = BuildTool
	itemtype.Use = UseBuildTool
	return AddItemType(itemtype)
}

func RegisterEnemyToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-enemy-tool"))
	itemtype.Name = "Dev Enemy Tool"
	itemtype.Use = UseEnemyTool
	return AddItemType(itemtype)
}

func RegisterPipeConnectToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-pipe-place-tool"))
	itemtype.Name = "Dev Pipe Connection Tool"
	itemtype.Use = UsePipeConnectionTool
	return AddItemType(itemtype)
}

func UseBuildTool(info ItemUseInfo) {
	didSelect := info.Inventory.PlacementGrid.TrySelect(info.CamCoords())
	if didSelect {
		return
	}
	didPlace := info.Inventory.PlacementGrid.TryPlace(info)
	_ = didPlace
}

func UseEnemyTool(info ItemUseInfo) {
	id := ui.EnemyToolActiveAgentID()
	world := info.WorldCoords()
	opts := agents.AgentCreationOptions{
		X: world.X(), Y: world.Y(),
	}
	info.World.Entities.Agents.Spawn(id, opts)
}

func UsePipeConnectionTool(info ItemUseInfo) {
	x32, y32 := cxmath.RoundVec2(info.WorldCoords())
	x := int(x32) ; y := int(y32)
	info.World.Planet.TryCyclePipeConnection(x,y)
}
