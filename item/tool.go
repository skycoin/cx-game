package item

import (
	"github.com/skycoin/cx-game/spriteloader"
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
