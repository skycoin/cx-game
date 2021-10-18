package item

import (
	"log"

	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/world"
)

func AddDrops() {
	regolithID, ok := world.IDFor("regolith")
	if !ok {
		log.Fatalf("cannot find tiletypeID for regolith")
	}
	world.AddDrop(regolithID, world.Drop{
		Count: 1, Probability: 0.5, Item: types.ItemTypeID(RockDustItemTypeID),
	})
}
