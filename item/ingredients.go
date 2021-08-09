package item

import (
	"github.com/skycoin/cx-game/engine/spriteloader"
)

func RegisterRockDustItemType() ItemTypeID {
	sprite := spriteloader.LoadSingleSprite(
		"./assets/item/rock1_dust.png", "rock-dust")
	itemtype := NewItemType(sprite)
	itemtype.Name = "Rock Dust"
	return AddItemType(itemtype)
}
