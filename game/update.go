package game

import (
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/sound"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/stars/starfield"
)

func Update(dt float32) {
	// player = findPlayer()
	timer.Accumulator += dt

	for timer.Accumulator >= constants.PHYSICS_TICK {
		FixedTick()
		timer.Accumulator -= constants.PHYSICS_TICK
	}
	components.Update(dt)
	ScreenManager.Update(dt)

	if Cam.IsFreeCam() {
		Cam.MoveCam(dt)
	} else {
		alpha := timer.GetTimeBetweenTicks() / constants.PHYSICS_TICK
		body :=
			World.Entities.Agents.FromID(playerAgentID).Transform

		var interpolatedPos cxmath.Vec2
		if !body.PrevPos.Equal(body.Pos) {
			interpolatedPos = body.PrevPos.Mult(1 - alpha).Add(body.Pos.Mult(alpha))

		} else {
			interpolatedPos = body.Pos
		}
		Cam.SetCameraPosition(interpolatedPos.X, interpolatedPos.Y)
	}

	World.Planet.Update(dt)

	Cam.Tick(dt)
	fps.Tick()
	ui.TickDialogueBoxes(dt)
	//obsolete
	particles.TickParticles(dt)

	sound.SetListenerPosition(player.GetAgent().Transform.Pos)
	//has to be after listener position is updated
	sound.Update()

	starfield.UpdateStarField(dt)
	ui.TickDamageIndicators(dt)
}

func FixedTick() {
	components.FixedUpdate()
	World.Planet.FixedUpdate()
	World.TimeState.Advance()
	ScreenManager.FixedUpdate()

	// physics.Simulate(constants.PHYSICS_TICK, &World.Planet)
	/*
		pickedUpItems := item.TickWorldItems(
			&World.Planet, physicsconstants.PHYSICS_TIMESTEP, player.Pos)
		for _, worldItem := range pickedUpItems {
			item.GetInventoryById(inventoryId).TryAddItem(worldItem.ItemTypeId)
		}
	*/
}
