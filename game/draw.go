package game

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/render/worldctx"
	"github.com/skycoin/cx-game/starfield"
	"github.com/skycoin/cx-game/world"
)

var debugTileInfo bool = true

func Draw() {
	// gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
	// gl.Disable(gl.BLEND)
	gl.ClearColor(7.0/255.0, 8.0/255.0, 25.0/255.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	render.FRAMEBUFFER_PLANET.Bind(gl.FRAMEBUFFER)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	render.FRAMEBUFFER_MAIN.Bind(gl.FRAMEBUFFER)
	gl.ClearColor(7.0/255.0, 8.0/255.0, 25.0/255.0, 1.0)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	render.SetCameraTransform(Cam.GetTransform())
	render.SetWorldWidth(float32(World.Planet.Width))
	// camCtx := baseCtx.PushView(Cam.GetView())
	worldCtx := worldctx.NewWorldRenderContext(Cam, &World.Planet)

	//draw starfield
	starfield.DrawStarField()

	//queue-draw world
	World.Planet.Draw(Cam, world.WindowLayer)
	World.Planet.Draw(Cam, world.BgLayer)
	World.Planet.Draw(Cam, world.PipeLayer)
	World.Planet.Draw(Cam, world.MidLayer)
	// draw lasers between mid and top layers.

	// particles.DrawTopParticles(camCtx)
	World.Planet.Draw(Cam, world.TopLayer)
	World.Planet.Draw(Cam, world.SuperLayer)

	item.DrawWorldItems(Cam)

	topLeftCtx :=
		render.CenterToTopLeft(win.DefaultRenderContext())
	ui.DrawString(
		fmt.Sprint(fps.CurFps),
		mgl32.Vec4{1, 0.2, 0.3, 1},
		ui.AlignCenter,
		topLeftCtx.PushLocal(mgl32.Translate3D(1, -5, 0)),
	)

	// FIXME: draw dialogue boxes uses alternate projection matrix;
	// restore original projection matrix

	//fix dialogboxdraw
	// ui.DrawDialogueBoxes(camCtx)
	Console.Draw(win.DefaultRenderContext())
	actualScreenSizeWidth, actualScreenSizeHeight := glfw.GetCurrentContext().GetFramebufferSize()
	physicalViewport := render.GetCurrentViewport()
	fmt.Println("physicalViewport: ", physicalViewport)
	virtualViewport :=
		render.Viewport{
			0, 0,
			int32(actualScreenSizeWidth),  //constants.VIRTUAL_VIEWPORT_WIDTH,
			int32(actualScreenSizeHeight), //constants.VIRTUAL_VIEWPORT_HEIGHT,
		}
	virtualViewport.Use()

	render.Flush(Cam.Zoom.Get())
	components.Draw_Queued(&World.Entities, Cam)
	//draw after flushing
	components.Draw(&World.Entities, Cam)
	gl.Disable(gl.BLEND)
	particles.DrawMidTopParticles(worldCtx)

	//draw lightmap
	World.Planet.DrawLighting(Cam, &World.TimeState)

	//post-process
	physicalViewport.Use()
	PostProcess()

	//draw ui
	DrawUI()

}

func PostProcess() {
	render.FRAMEBUFFER_SCREEN.Bind(gl.FRAMEBUFFER)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	render.ScreenShader.Use()
	render.ScreenShader.SetInt("u_screenTexture", 0)

	render.SetupShockwave()

	// gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, render.RENDERTEXTURE_MAIN)
	gl.BindVertexArray(render.Quad2Vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func DrawUI() {
	baseCtx := win.DefaultRenderContext()
	topLeftCtx :=
		render.CenterToTopLeft(baseCtx)
	if debugTileInfo {
		ui.DrawString(
			tileText,
			mgl32.Vec4{0.3, 0.9, 0.4, 1},
			ui.AlignCenter,
			topLeftCtx.PushLocal(mgl32.Translate3D(25, -5, 0)),
		)
	}
	ui.DrawAgentHUD(player)
	inventory := item.GetInventoryById(player.InventoryID)
	invCameraTransform := Cam.GetTransform().Inv()
	inventory.Draw(baseCtx, invCameraTransform)
	ui.DrawDamageIndicators(invCameraTransform)
	render.FlushUI()
	glfw.PollEvents()
	win.Window.SwapBuffers()
}
