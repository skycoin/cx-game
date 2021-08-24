package game

import (
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/sound"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/engine/ui/console"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/stars/starfield"
)

func Update(dt float32) {
	player = findPlayer()
	FixedUpdate(dt)
	// physics.Simulate(dt, CurrentPlanet)
	components.Update(dt)
	if Cam.IsFreeCam() {
		//player.Controlled = false
		Cam.MoveCam(dt)
	} else {
		//player.Controlled = true
		//playerPos := player.InterpolatedTransform.Col(3).Vec2()
		alpha := timer.GetTimeBetweenTicks() / constants.PHYSICS_TICK
		body :=
			World.Entities.Agents.FromID(playerAgentID).PhysicsState
		interpolatedPos := body.PrevPos.Mult(1 - alpha).Add(body.Pos.Mult(alpha))

		Cam.SetCameraPosition(interpolatedPos.X, interpolatedPos.Y)
	}
	World.Planet.Update(dt)
	Cam.Tick(dt)
	fps.Tick()
	ui.TickDialogueBoxes(dt)
	particles.TickParticles(dt)

	sound.SetListenerPosition(player.PhysicsState.Pos)
	//has to be after listener position is updated
	sound.Update()

	commandContext := console.NewCommandContext()
	commandContext.World = &World
	commandContext.Player = player
	Console.Update(window, commandContext)

	starfield.UpdateStarField(dt)

	mousePos := input.GetMousePos()
	item.GetInventoryById(player.InventoryID).PlacementGrid.UpdatePreview(
		&World,
		mousePos.X(),
		mousePos.Y(),
		Cam,
	)

}
