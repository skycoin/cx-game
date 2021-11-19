package game

// import (
// 	"github.com/go-gl/glfw/v3.3/glfw"
// 	"github.com/skycoin/cx-game/engine"
// 	"github.com/skycoin/cx-game/engine/camera"
// 	"github.com/skycoin/cx-game/engine/input"
// 	"github.com/skycoin/cx-game/engine/sound"
// 	"github.com/skycoin/cx-game/engine/ui"
// 	"github.com/skycoin/cx-game/item"
// 	"github.com/skycoin/cx-game/render"
// 	"github.com/skycoin/cx-game/world"
// )

// func ProcessInput() {
// 	ScreenManager.ProcessInput()
// }

// 	switch input.GetInputContext() {
// 	case input.GAME:
// 		if input.GetButtonDown("mute") {
// 			sound.ToggleMute()
// 		}
// 		if input.GetButtonDown("freecam-on") {
// 			Cam.TurnOnFreeCam()
// 		}
// 		if input.GetButtonDown("inventory-grid") {
// 			inventory := item.GetInventoryById(player.InventoryID)
// 			inventory.IsOpen = !inventory.IsOpen
// 		}
// 		inventory := item.GetInventoryById(player.InventoryID)
// 		inventory.TrySelectSlot(input.GetLastKey())
// 		if input.GetButtonDown("enemy-tool-scroll-down") {
// 			ui.EnemyToolScrollDown()
// 			inventory.TryScrollDown()
// 		}
// 		if input.GetButtonDown("enemy-tool-scroll-up") {
// 			ui.EnemyToolScrollUp()
// 			inventory.TryScrollUp()
// 		}
// 		if input.GetButtonDown("toggle-zoom") {
// 			Cam.CycleZoom()
// 		}

// 		if input.GetButtonDown("toggle-texture-filtering") {
// 			render.ToggleFiltering()
// 		}
// 		if input.GetButtonDown("toggle-bbox") {
// 			render.ToggleBBox()
// 		}
// 		if input.GetButtonDown("cycle-pixel-snap") {
// 			render.CyclePixelSnap()
// 		}
// 		if input.GetButtonDown("cycle-camera-snap") {
// 			camera.CycleSnap()
// 		}
// 		if input.GetButtonDown("switch-skylight") {
// 			world.SwitchNeighbourCount(&World.Planet)
// 		}

// 		if input.GetButtonDown("set-camera-player") {
// 			Cam.SwitchToPlayer()
// 		}
// 		if input.GetButtonDown("set-camera-target") {
// 			Cam.SwitchToTarget()
// 		}
// 		if input.GetKeyDown(glfw.KeyN) {
// 			world.ToggleSmoothLighting()
// 		}

// 		if input.GetKeyDown(glfw.KeyB) {
// 			Cam.Shake()
// 		}
// 		if input.GetKeyDown(glfw.KeyV) {
// 			Cam.ShockwaveVec(player.Transform.Pos)
// 			// Cam.ScreenPos(player.PhysicsState.Pos.Mgl32())
// 		}

// 	case input.FREECAM:
// 		if input.GetButtonDown("freecam-off") {
// 			Cam.TurnOffFreeCam()
// 		}
// 		if input.GetButtonDown("toggle-bbox") {
// 			render.ToggleBBox()
// 		}
// 		if input.GetButtonDown("cycle-pixel-snap") {
// 			render.CyclePixelSnap()
// 		}
// 		if input.GetButtonDown("cycle-camera-snap") {
// 			camera.CameraSnapped = !camera.CameraSnapped
// 			camera.CycleSnap()
// 		}
// 	}

// 	if input.GetButtonDown("toggle-log") {
// 		engine.ToggleLogging()
// 	}

// 	if input.GetKeyDown(glfw.KeyHome) {
// 		debugTileInfo = !debugTileInfo
// 	}

// }
