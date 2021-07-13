package game

import (
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/enemies"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/sound"
	"github.com/skycoin/cx-game/starfield"
	"github.com/skycoin/cx-game/ui"
)

func Update(dt float32) {
	player.Update(dt, CurrentPlanet)
	physics.Simulate(dt, CurrentPlanet)
	components.Update(CurrentPlanet.WorldState,player)
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
	pickedUpItems := item.TickWorldItems(CurrentPlanet, dt, player.Pos)
	for _, worldItem := range pickedUpItems {
		item.GetInventoryById(inventoryId).TryAddItem(worldItem.ItemTypeId)
	}
	enemies.TickBasicEnemies(CurrentPlanet, dt, player, catIsScratching)

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
