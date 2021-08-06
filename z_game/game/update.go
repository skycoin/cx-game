package game

import (
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/sound"
	"github.com/skycoin/cx-game/starfield"
	"github.com/skycoin/cx-game/ui"
	"github.com/skycoin/cx-game/ui/console"
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
		playerPos :=
			World.Entities.Agents.FromID(playerAgentID).PhysicsState.Pos
		Cam.SetCameraPosition(playerPos.X, playerPos.Y)
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
	Console.Update(window, commandContext)

	starfield.UpdateStarField(dt)

}

// type mouseDraws struct {
// 	xpos float32
// 	ypos float32
// }
