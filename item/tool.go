package item

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/render"
)

func RegisterFurnitureToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-furniture-tool-2"))
	itemtype.Name = "Dev Furniture Tool"
	itemtype.Category = BuildTool
	itemtype.Use = UseFurnitureTool
	return AddItemType(itemtype)
}

func RegisterTileToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-tile-tool"))
	itemtype.Name = "Dev Tile Tool"
	itemtype.Category = BuildTool
	itemtype.Use = UseTileTool
	return AddItemType(itemtype)

}

func RegisterEnemyToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-enemy-tool"))
	itemtype.Name = "Dev Enemy Tool"
	itemtype.Use = UseEnemyTool
	return AddItemType(itemtype)
}

func UseFurnitureTool(info ItemUseInfo) {
	didSelect := info.Inventory.PlacementGrid.TrySelect(info.CamCoords())
	if didSelect {
		return
	}
	didPlace := info.Inventory.PlacementGrid.TryPlace(info)
	_ = didPlace
}

func UseTileTool(info ItemUseInfo) {
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
