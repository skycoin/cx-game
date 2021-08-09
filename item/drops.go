package item

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/world"
)

func AddDrops() {
	world.AddDrop(world.IDFor("regolith"), world.Drop{
		Count: 1, Probability: 0.5, Item: types.ItemTypeID(RockDustItemTypeID),
	})
}
