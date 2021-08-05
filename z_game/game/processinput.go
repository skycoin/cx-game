package game

import (
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/skycoin/cx-game/input"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/sound"
	"github.com/skycoin/cx-game/starfield"
	"github.com/skycoin/cx-game/ui"
)

func ProcessInput() {
	if Console.IsActive() { return }
	/*
	if input.GetButtonDown("switch-helmet") {
		player.SetHelmNext()
	}
	if input.GetButtonDown("switch-suit") {
		player.SetSuitNext()
	}
	*/
	/*
	if input.GetButtonDown("jump") {
		didJump := player.Jump()
		if didJump {
			ui.PlaceDialogueBox(
				"*jump*", ui.AlignRight, 1,
				mgl32.Translate3D(
					player.Pos.X,
					player.Pos.Y,
					0,
				),
			)
			sound.PlaySound("player_jump", sound.SoundOptions{Pitch: 1.5})
		}
	}
	if input.GetButtonDown("fly") {
		player.ToggleFlying()
	}
	*/
	if input.GetButtonDown("mute") {
		sound.ToggleMute()
	}
	if input.GetButtonDown("freecam") {
		Cam.ToggleFreeCam()
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
	if input.GetButtonDown("enemy-tool-scroll-down") {
		ui.EnemyToolScrollDown()
	}
	if input.GetButtonDown("enemy-tool-scroll-up") {
		ui.EnemyToolScrollUp()
	}
	inventory := item.GetInventoryById(player.InventoryID)
	inventory.TrySelectSlot(input.GetLastKey())
}
