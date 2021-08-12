package game

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/sound"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/starfield"
)

func ProcessInput() {
	if Console.IsActive() {
		return
	}
	if input.GetButtonDown("mute") {
		sound.ToggleMute()
	}
	if input.GetButtonDown("freecam") {
		Cam.ToggleFreeCam()
		if input.GetInputContext() == input.GAME{
			input.SetInputContext(input.FREECAM)
		}else if input.GetInputContext() == input.FREECAM{
			input.SetInputContext(input.GAME)
		}
	}
	if input.GetButtonDown("inventory-grid") {
		inventory := item.GetInventoryById(player.InventoryID)
		inventory.IsOpen = !inventory.IsOpen
	}
	if input.GetKeyDown(glfw.KeyL) {
		starfield.SwitchBackgrounds(starfield.BACKGROUND_NEBULA)
	}
	if input.GetKeyDown(glfw.KeyO) {
		starfield.SwitchBackgrounds(starfield.BACKGROUND_VOID)
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
}
