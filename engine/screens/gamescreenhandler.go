package screens

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/components"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/input/inputhandler"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/particles"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/render/worldctx"
	"github.com/skycoin/cx-game/stars/starfield"
	"github.com/skycoin/cx-game/world"
)

const(
	
)
type GameScreenHandler struct {
	inputHandler *inputhandler.AgentInputHandler
}

func NewGameScreenHandler() *GameScreenHandler {
	newScreenHandler := GameScreenHandler{
		inputHandler: inputhandler.NewAgentInputHandler(),
	}

	return &newScreenHandler
}

func (g *GameScreenHandler) RegisterAction(info inputhandler.ActionInfo) {
	g.inputHandler.RegisterAction(info)
}

func (g *GameScreenHandler) ProcessInput() {
	g.inputHandler.ProcessEvents()
	g.inputHandler.UpdateKeyState()
}
func (g *GameScreenHandler) FixedUpdate() {
	//player.FixedTick(&World.Planet)
	components.FixedUpdate()

	g.player.SetControlState(g.inputHandler.AgentControlState)

	// physics.Simulate(constants.PHYSICS_TICK, &World.Planet)
	/*
		pickedUpItems := item.TickWorldItems(
			&World.Planet, physicsconstants.PHYSICS_TIMESTEP, player.Pos)
		for _, worldItem := range pickedUpItems {
			item.GetInventoryById(inventoryId).TryAddItem(worldItem.ItemTypeId)
		}
	*/
}

func (g *GameScreenHandler) Update(dt float32) {
	timer.Accumulator += dt

	for timer.Accumulator >= constants.PHYSICS_TICK {

		g.FixedUpdate()
		timer.Accumulator -= constants.PHYSICS_TICK
	}
	components.Update(dt)

	if g.Cam.IsFreeCam() {
		//player.Controlled = false
		g.Cam.MoveCam(dt)
	} else {
		//player.Controlled = true
		//playerPos := player.InterpolatedTransform.Col(3).Vec2()
		alpha := timer.GetTimeBetweenTicks() / constants.PHYSICS_TICK
		body :=
			g.World.Entities.Agents.FromID(g.player.GetAgent().AgentId).Transform

		var interpolatedPos cxmath.Vec2
		if !body.PrevPos.Equal(body.Pos) {
			interpolatedPos = body.PrevPos.Mult(1 - alpha).Add(body.Pos.Mult(alpha))

		} else {
			interpolatedPos = body.Pos
		}
		g.Cam.SetCameraPosition(interpolatedPos.X, interpolatedPos.Y)
	}
}

func (g *GameScreenHandler) RenderWorld() {

}

func (g *GameScreenHandler) RenderUI() {

}

func (g *GameScreenHandler) Render() {
	g.RenderWorld()
	g.RenderUI()
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
	// g.CamCtx := baseCtx.PushView(g.g.Cam.GetView())
	worldCtx := worldctx.NewWorldRenderContext(Cam, &World.Planet)

	//draw starfield
	starfield.DrawStarField()

	//queue-draw world
	World.Planet.Draw(Cam, world.WindowLayer)
	World.Planet.Draw(Cam, world.BgLayer)
	World.Planet.Draw(Cam, world.PipeLayer)
	World.Planet.Draw(Cam, world.MidLayer)

	// draw lasers between mid and top layers.

	// particles.DrawTopParticles(g.CamCtx)
	World.Planet.Draw(Cam, world.TopLayer)

	item.DrawWorldItems(Cam)

	// FIXME: draw dialogue boxes uses alternate projection matrix;
	// restore original projection matrix

	//fix dialogboxdraw
	// ui.DrawDialogueBoxes(g.CamCtx)

	Console.Draw(win.DefaultRenderContext())

	physicalViewport := render.GetCurrentViewport()
	virtualViewport :=
		render.Viewport{
			0, 0,
			constants.VIRTUAL_VIEWPORT_WIDTH,
			constants.VIRTUAL_VIEWPORT_HEIGHT,
		}
	virtualViewport.Use()
	components.Draw_Queued(&World.Entities, Cam)
	render.Flush(Cam.Zoom.Get())

	// tile := g.World.Planet.GetTile(35, 4, world.TopLayer)

	// fmt.Println(tile.Name)
	//draw after flushing
	components.Draw(World.Entities, Cam)
	gl.Disable(gl.BLEND)
	particles.DrawMidTopParticles(worldCtx)

	//draw lightmap
	World.Planet.DrawLighting(Cam, &World.TimeState)

	//post-process
	physicalViewport.Use()
	g.PostProcess()
	//draw ui
}

func (g *GameScreenHandler) PostProcess() {
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

var debugTileInfo bool = false

func (g *GameScreenHandler) DrawUI() {
	baseCtx := g.win.DefaultRenderContext()
	topLeftCtx :=
		render.CenterToTopLeft(baseCtx)

	//todo pass fps instance here

	ui.DrawString(
		fmt.Sprint(g.fps.CurFps),
		mgl32.Vec4{1, 0.2, 0.3, 1},
		ui.AlignCenter,
		topLeftCtx.PushLocal(mgl32.Translate3D(1, -5, 0)),
	)
	if debugTileInfo {
		ui.DrawString(
			// tileText,
			"somtn",
			mgl32.Vec4{0.3, 0.9, 0.4, 1},
			ui.AlignCenter,
			topLeftCtx.PushLocal(mgl32.Translate3D(25, -5, 0)),
		)
	}
	ui.DrawAgentHUD(g.player.GetAgent())
	inventory := item.GetInventoryById(g.player.GetAgent().InventoryID)
	invCameraTransform := g.Cam.GetTransform().Inv()
	inventory.Draw(baseCtx, invCameraTransform)
	ui.DrawDamageIndicators(invCameraTransform)
	render.FlushUI()
}
