package item

import (
	"github.com/skycoin/cx-game/render"
)

func RegisterRockDustItemType() ItemTypeID {
	itemtype := NewItemType(render.GetSpriteIDByName("rock-dust"))
	itemtype.Name = "Rock Dust"
	return AddItemType(itemtype)
}
