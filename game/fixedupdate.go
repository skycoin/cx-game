package game

import (
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/constants"

	//"github.com/skycoin/cx-game/item"

	"github.com/skycoin/cx-game/physics/timer"
)

func FixedUpdate(dt float32) {
	timer.Accumulator += dt

	for timer.Accumulator >= constants.PHYSICS_TICK {

		FixedTick()
		timer.Accumulator -= constants.PHYSICS_TICK
	}
}

func FixedTick() {
	//player.FixedTick(&World.Planet)
	components.FixedUpdate()
	World.Planet.UpdateSkyLight(1000)
	// physics.Simulate(constants.PHYSICS_TICK, &World.Planet)
	/*
		pickedUpItems := item.TickWorldItems(
			&World.Planet, physicsconstants.PHYSICS_TIMESTEP, player.Pos)
		for _, worldItem := range pickedUpItems {
			item.GetInventoryById(inventoryId).TryAddItem(worldItem.ItemTypeId)
		}
	*/
}
