package game

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/input"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/world"
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

func mouseReleaseCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	screenX := float32(input.GetMouseX()-float64(win.Width)/2) / Cam.Zoom // adjust mouse position with zoom
	screenY := (float32(input.GetMouseY()-float64(win.Height)/2) * -1) / Cam.Zoom

	if isInventoryGridVisible {
		inventory := item.GetInventoryById(inventoryId)
		inventory.TryMoveSlot(screenX, screenY, Cam, CurrentPlanet, player)
	}
}

func mousePressCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	// we only care about mousedown events for now
	if a != glfw.Press {
		return
	}

	screenX := float32(((input.GetMouseX()-float64(widthOffset))/float64(scale) - float64(win.Width)/2)) / Cam.Zoom // adjust mouse position with zoom
	screenY := float32(((input.GetMouseY()-float64(heightOffset))/float64(scale)-float64(win.Height)/2)*-1) / Cam.Zoom

	didSelectPaleteTile := tilePaletteSelector.TrySelectTile(screenX, screenY)
	if didSelectPaleteTile {
		return
	}

	if tilePaletteSelector.IsMultiTileSelected() {
		didPlaceMultiTile := CurrentPlanet.TryPlaceMultiTile(
			screenX, screenY,
			world.Layer(tilePaletteSelector.LayerIndex),
			tilePaletteSelector.GetSelectedMultiTile(),
			Cam,
		)
		if didPlaceMultiTile {
			return
		}
	} else {
		didPlaceTile := CurrentPlanet.TryPlaceTile(
			screenX, screenY,
			world.Layer(tilePaletteSelector.LayerIndex),
			tilePaletteSelector.GetSelectedTile(),
			Cam,
		)
		if didPlaceTile {
			return
		}
	}

	if isInventoryGridVisible {
		inventory := item.GetInventoryById(inventoryId)
		clickedSlot :=
			inventory.TryClickSlot(screenX, screenY, Cam, CurrentPlanet, player)
		if clickedSlot {
			return
		}
	}

	item.GetInventoryById(inventoryId).
		TryUseItem(screenX, screenY, Cam, CurrentPlanet, player)
}

var (
	widthOffset, heightOffset int32
	scale                     float32 = 1
)

func windowSizeCallback(window *glfw.Window, width, height int) {

	// gl.Viewport(0, 0, int32(width), int32(height))
	scaleToFitWidth := float32(width) / float32(win.Width)
	scaleToFitHeight := float32(height) / float32(win.Height)
	scale = cxmath.Min(scaleToFitHeight, scaleToFitWidth)

	widthOffset = int32((float32(width) - float32(win.Width)*scale) / 2)
	heightOffset = int32((float32(height) - float32(win.Height)*scale) / 2)
	//correct mouse offsets
	input.UpdateMouseCoords(widthOffset, heightOffset, scale)

	gl.Viewport(widthOffset, heightOffset, int32(float32(win.Width)*scale), int32(float32(win.Height)*scale))
	// win.Width = width
	// win.Height = height
}

func scrollCallback(w *glfw.Window, xOff, yOff float64) {
	Cam.SetCameraZoomPosition(float32(yOff))
}