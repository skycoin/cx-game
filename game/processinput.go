package game

import (
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/sound"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/render"
)

func ProcessInput() {
	switch input.GetInputContext() {
	case input.GAME:
		if input.GetButtonDown("mute") {
			sound.ToggleMute()
		}
		if input.GetButtonDown("freecam-on") {
			Cam.TurnOnFreeCam()
		}
		if input.GetButtonDown("inventory-grid") {
			inventory := item.GetInventoryById(player.InventoryID)
			inventory.IsOpen = !inventory.IsOpen
		}
		inventory := item.GetInventoryById(player.InventoryID)
		inventory.TrySelectSlot(input.GetLastKey())
		if input.GetButtonDown("enemy-tool-scroll-down") {
			ui.EnemyToolScrollDown()
			inventory.TryScrollDown()
		}
		if input.GetButtonDown("enemy-tool-scroll-up") {
			ui.EnemyToolScrollUp()
			inventory.TryScrollUp()
		}
		if input.GetButtonDown("toggle-zoom") {
			Cam.CycleZoom()
		}

		if input.GetButtonDown("toggle-texture-filtering") {
			render.ToggleFiltering()
		}
		if input.GetButtonDown("toggle-bbox") {
			render.ToggleBBox()
		}
	case input.FREECAM:
		if input.GetButtonDown("freecam-off") {
			Cam.TurnOffFreeCam()
		}
		if input.GetButtonDown("toggle-bbox") {
			render.ToggleBBox()
		}
	}

	input.Reset()

}
