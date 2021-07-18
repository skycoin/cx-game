package game

import (
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/sound"
	"github.com/skycoin/cx-game/starfield"
	"github.com/skycoin/cx-game/ui"
)

func Update(dt float32) {
	FixedUpdate(dt)
	// physics.Simulate(dt, CurrentPlanet)
	components.Update(dt)
	if Cam.IsFreeCam() {
		player.Controlled = false
		Cam.MoveCam(dt)
	} else {
		player.Controlled = true
		playerPos := player.InterpolatedTransform.Col(3).Vec2()
		Cam.SetCameraPosition(playerPos.X(), playerPos.Y())
	}
	Cam.Tick(dt)
	fps.Tick()
	ui.TickDialogueBoxes(dt)
	particles.TickParticles(dt)

	sound.SetListenerPosition(player.Pos)
	//has to be after listener position is updated
	sound.Update()

	starfield.UpdateStarField(dt)
	catIsScratching = false

}

// type mouseDraws struct {
// 	xpos float32
// 	ypos float32
// }
