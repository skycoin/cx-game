package item

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/ui"
)

func UseBuildTool(info ItemUseInfo) {
	didSelect := info.Inventory.PlacementGrid.TrySelect(info.CamCoords())
	if didSelect { return }
	didPlace := info.Inventory.PlacementGrid.TryPlace(info)
	_ = didPlace
}

func RegisterBuildToolItemType() ItemTypeID {
	sprite := spriteloader.LoadSingleSprite(
		"./assets/item/dev-build.png", "dev-build-tool" )
	itemtype := NewItemType(sprite)
	itemtype.Name = "Dev Build Tool"
	itemtype.Category = BuildTool
	itemtype.Use = UseBuildTool
	return AddItemType(itemtype)
}

func UseEnemyTool(info ItemUseInfo) {
	id := ui.EnemyToolActiveAgentID()
	world := info.WorldCoords()
	opts := agents.AgentCreationOptions {
		X: world.X(), Y: world.Y(),
	}
	info.Planet.WorldState.AgentList.Spawn(id, opts)
}

func RegisterEnemyToolItemType() ItemTypeID {
	sprite := spriteloader.LoadSingleSprite(
		"./assets/item/dev-enemy.png", "dev-enemy-tool" )
	itemtype := NewItemType(sprite)
	itemtype.Name = "Dev Enemy Tool"
	itemtype.Use = UseEnemyTool
	return AddItemType(itemtype)
}
