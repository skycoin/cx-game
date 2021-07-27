package game

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/render/worldctx"
	"github.com/skycoin/cx-game/starfield"
	"github.com/skycoin/cx-game/ui"
	"github.com/skycoin/cx-game/world"
)

func Draw() {
	gl.ClearColor(7.0/255.0, 8.0/255.0, 25.0/255.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	baseCtx := win.DefaultRenderContext()
	baseCtx.Projection = Cam.GetProjectionMatrix()
	camCtx := baseCtx.PushView(Cam.GetView())
	worldCtx := worldctx.NewWorldRenderContext(Cam, CurrentPlanet)

	starfield.DrawStarField()
	CurrentPlanet.Draw(Cam, world.BgLayer)
	CurrentPlanet.Draw(Cam, world.MidLayer)
	// draw lasers between mid and top layers.
	particles.DrawMidTopParticles(worldCtx)
	CurrentPlanet.Draw(Cam, world.TopLayer)
	particles.DrawTopParticles(camCtx)

	item.DrawWorldItems(Cam)
	components.Draw(CurrentPlanet.WorldState, Cam)
	player.Draw(Cam, CurrentPlanet)
	ui.DrawHUD(player.GetHUDState())

	ui.DrawString(
		fmt.Sprint(fps.CurFps),
		mgl32.Vec4{1, 0.2, 0.3, 1},
		ui.AlignCenter,
		win.DefaultRenderContext().PushLocal(mgl32.Translate3D(-11.5, 5, 0)),
	)

	// tile - air line (green)
	collidingTileLines := CurrentPlanet.GetCollidingTilesLinesRelative(
		int(player.Pos.X), int(player.Pos.Y))
	if len(collidingTileLines) > 2 {
		Cam.DrawLines(collidingTileLines, []float32{0.0, 1.0, 0.0}, baseCtx)
	}

	// body bounding box (blue)
	Cam.DrawLines(player.GetBBoxLines(), []float32{0.0, 0.0, 1.0}, baseCtx)

	// colliding line from body (red)
	collidingLines := player.GetCollidingLines()
	if len(collidingLines) > 2 {
		Cam.DrawLines(collidingLines, []float32{1.0, 0.0, 0.0}, baseCtx)
	}

	ui.DrawDialogueBoxes(camCtx)
	// FIXME: draw dialogue boxes uses alternate projection matrix;
	// restore original projection matrix

	inventory := item.GetInventoryById(inventoryId)
	inventory.Draw(win.DefaultRenderContext())
	tilePaletteSelector.Draw(win.DefaultRenderContext())

	glfw.PollEvents()
	win.Window.SwapBuffers()
}
