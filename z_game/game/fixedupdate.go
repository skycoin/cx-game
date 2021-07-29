package game

import (
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/constants/physicsconstants"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/physics/timer"
)

func FixedUpdate(dt float32) {
	timer.Accumulator += dt

	for timer.Accumulator >= physicsconstants.PHYSICS_TIMESTEP {

		FixedTick()
		timer.Accumulator -= physicsconstants.PHYSICS_TIMESTEP
	}
}

func FixedTick() {
	player.FixedTick(&World.Planet)
	components.FixedUpdate()
	physics.Simulate(physicsconstants.PHYSICS_TIMESTEP, &World.Planet)
	pickedUpItems := item.TickWorldItems(
		&World.Planet, physicsconstants.PHYSICS_TIMESTEP, player.Pos)
	for _, worldItem := range pickedUpItems {
		item.GetInventoryById(inventoryId).TryAddItem(worldItem.ItemTypeId)
	}
}
