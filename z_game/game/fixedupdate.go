package game

import (
	"github.com/skycoin/cx-game/constants/physicsconstants"
	"github.com/skycoin/cx-game/enemies"
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
	player.FixedTick(CurrentPlanet)
	physics.Simulate(physicsconstants.PHYSICS_TIMESTEP, CurrentPlanet)
	enemies.TickBasicEnemies(CurrentPlanet, physicsconstants.PHYSICS_TIMESTEP, player, catIsScratching)
	pickedUpItems := item.TickWorldItems(CurrentPlanet, physicsconstants.PHYSICS_TIMESTEP, player.Pos)
	for _, worldItem := range pickedUpItems {
		item.GetInventoryById(inventoryId).TryAddItem(worldItem.ItemTypeId)
	}
}
