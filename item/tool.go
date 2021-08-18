package item

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/render"
)

func UseBuildTool(info ItemUseInfo) {
	didSelect := info.Inventory.PlacementGrid.TrySelect(info.CamCoords())
	if didSelect {
		return
	}
	didPlace := info.Inventory.PlacementGrid.TryPlace(info)
	_ = didPlace
}

func RegisterBuildToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-build-furniture-tool-2"))
	itemtype.Name = "Dev Build Tool"
	itemtype.Category = BuildTool
	itemtype.Use = UseBuildTool
	return AddItemType(itemtype)
}

func UseEnemyTool(info ItemUseInfo) {
	id := ui.EnemyToolActiveAgentID()
	world := info.WorldCoords()
	opts := agents.AgentCreationOptions{
		X: world.X(), Y: world.Y(),
	}
	info.World.Entities.Agents.Spawn(id, opts)
}

func RegisterEnemyToolItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("dev-enemy-tool"))
	itemtype.Name = "Dev Enemy Tool"
	itemtype.Use = UseEnemyTool
	return AddItemType(itemtype)
}
