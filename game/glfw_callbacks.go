package game

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/world"
)

var leftMouseDown = false

func mouseButtonCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	if a == glfw.Press {
		mousePressCallback(w, b, a, mk)
		leftMouseDown = true
	}
	if a == glfw.Release {
		mouseReleaseCallback(w, b, a, mk)
		leftMouseDown = false
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
	input.MousePressed = false
	mousePos := input.GetMousePos()

	inventory := item.GetInventoryById(player.InventoryID)
	player := findPlayer()
	inventory.OnReleaseMouse(mousePos.X(), mousePos.Y(), Cam, &World.Planet, player)
}

func mousePressCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	input.MousePressed = true

	// we only care about mousedown events for now
	if a != glfw.Press {
		return
	}

	mousePos := input.GetMousePos()

	// return
	inventory := item.GetInventoryById(player.InventoryID)

	// only if dev destroy tool selected
	if inventory.SelectedBarSlotIndex == 0 {
		worldCoords := Cam.GetTransform().Mul4x1(mousePos.Mul(1.0/32).Vec4(0, 1))

		worldX, worldY := cxmath.RoundVec2(worldCoords.Vec2())

		tile := World.Planet.GetTile(int(worldX), int(worldY), world.TopLayer)
		if tile.Name == "" {
			tile = World.Planet.GetTile(int(worldX), int(worldY), world.MidLayer)
			if tile.Name == "" {
				tile = World.Planet.GetTile(int(worldX), int(worldY), world.BgLayer)
				if tile.Name != "" {
					item.SelectedLayer = world.BgLayer
				}
			} else {
				item.SelectedLayer = world.MidLayer
			}
		} else {
			item.SelectedLayer = world.TopLayer
		}
	}

	//for dev destroy tool

	if b == glfw.MouseButtonRight {
		inventory.TryCancelSelect()
		return
	}
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

func cursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	input.MouseCoords.X = float32(xpos)
	input.MouseCoords.Y = float32(ypos)

	mousePos := input.GetMousePos()

	if leftMouseDown {
		item.GetInventoryById(player.InventoryID).
			TryUseBuildItem(mousePos.X(), mousePos.Y(), Cam, &World, player)
	}

	worldCoords := Cam.GetTransform().Mul4x1(mousePos.Mul(1.0/32).Vec4(0, 1)).Vec2()

	worldX, worldY := cxmath.RoundVec2(worldCoords)
	// tile := World.Planet.GetTile(int(worldCoords[0]), int(worldCoords[1]), world.TopLayer)
	tile := World.Planet.GetTile(int(worldX), int(worldY), world.TopLayer)
	idx := World.Planet.GetTileIndex(int(worldX), int(worldY))
	if tile == nil {
		return
	}
	tileText = fmt.Sprint(tile.TileCollisionType, "   ", tile.Name, "    ", World.Planet.LightingValues[idx].GetEnvLight(), "    ", World.Planet.LightingValues[idx].GetSkyLight(), "  |  ", worldX, "  ", worldY)

}

var tileText string

func windowSizeCallback(window *glfw.Window, width, height int) {
	// "physical" dimensions describe actual window size
	// "virtual" dimensions describe scaling of both world and UI
	// physical determines resolution.
	// virtual determines how big things are.
	// physicalWidth := float32(width)
	// physicalHeight := float32(height)
	// virtualWidth := float32(win.Width)
	// virtualHeight := float32(win.Height)

	// scaleToFitWidth := physicalWidth / virtualWidth
	// scaleToFitHeight := physicalHeight / virtualHeight
	// // scale to fit entire virtual window in physical window
	// scale := cxmath.Min(scaleToFitHeight, scaleToFitWidth)

	// // scale up virtual dimensions to fit in physical dimensions.
	// // in case of aspect ratio mismatch, black bars will appear
	// viewportWidth := int32(virtualWidth * scale)
	// viewportHeight := int32(virtualHeight * scale)

	// // store offsets for transitioning from physical to virtual mouse coords
	// // TODO store virtual coords, NOT physical coords
	// widthOffset := (int32(physicalWidth) - viewportWidth) / 2
	// heightOffset := (int32(physicalHeight) - viewportHeight) / 2

	// widthOffset += heightOffset
	// gl.Viewport(widthOffset, heightOffset, viewportWidth, viewportHeight)
}

func scrollCallback(w *glfw.Window, xOff, yOff float64) {
	Cam.SetCameraZoomPosition(float32(yOff))
}
