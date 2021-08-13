package game

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/item"
)

func mouseButtonCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	if a == glfw.Press {
		mousePressCallback(w, b, a, mk)
	}
	if a == glfw.Release {
		mouseReleaseCallback(w, b, a, mk)
	}
}

/*
// mouse position relative to screen
func screenPos() (float32, float32) {
	w := float32(win.Width)
	h := float32(win.Height)


	// adjust mouse position with zoom
	screenX :=
		((input.GetMouseX()-float32(widthOffset))/scale - w/2) / Cam.Zoom
	screenY :=
		-((input.GetMouseY()-float32(heightOffset))/scale - h/2) / Cam.Zoom
	return screenX, screenY
}
*/

func mouseReleaseCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	mousePos := input.GetMousePos()

	inventory := item.GetInventoryById(player.InventoryID)
	player := findPlayer()
	inventory.OnReleaseMouse(mousePos.X(), mousePos.Y(), Cam, &World.Planet, player)
}

func mousePressCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	// we only care about mousedown events for now
	if a != glfw.Press {
		return
	}

	mousePos := input.GetMousePos()

	inventory := item.GetInventoryById(player.InventoryID)
	clickedSlot :=
		inventory.TryClickSlot(
			mousePos.X(), mousePos.Y(), Cam, &World.Planet, player,
		)
	if clickedSlot {
		return
	}

	player := World.Entities.Agents.FromID(playerAgentID)
	item.GetInventoryById(player.InventoryID).
		TryUseItem(mousePos.X(), mousePos.Y(), Cam, &World, player)
}

var (
	// what actually are these???
	widthOffset, heightOffset int32
	scale                     float32 = 1
)

func windowSizeCallback(window *glfw.Window, width, height int) {
	// "physical" dimensions describe actual window size
	// "virtual" dimensions describe scaling of both world and UI
	// physical determines resolution.
	// virtual determines how big things are.
	physicalWidth := float32(width)
	physicalHeight := float32(height)
	virtualWidth := float32(win.Width)
	virtualHeight := float32(win.Height)

	scaleToFitWidth := physicalWidth / virtualWidth
	scaleToFitHeight := physicalHeight / virtualHeight
	// scale to fit entire virtual window in physical window
	scale = cxmath.Min(scaleToFitHeight, scaleToFitWidth)

	// scale up virtual dimensions to fit in physical dimensions.
	// in case of aspect ratio mismatch, black bars will appear
	viewportWidth := int32(virtualWidth * scale)
	viewportHeight := int32(virtualHeight * scale)

	// store offsets for transitioning from physical to virtual mouse coords
	// TODO store virtual coords, NOT physical coords
	widthOffset = (int32(physicalWidth) - viewportWidth) / 2
	heightOffset = (int32(physicalHeight) - viewportHeight) / 2

	gl.Viewport(widthOffset, heightOffset, viewportWidth, viewportHeight)
}

func scrollCallback(w *glfw.Window, xOff, yOff float64) {
	Cam.SetCameraZoomPosition(float32(-yOff))
	input.SetCamZoom(Cam.Zoom)
}
