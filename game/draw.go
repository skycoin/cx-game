package game

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/item"
	"github.com/skycoin/cx-game/render"
)

var debugTileInfo bool = true

func Draw() {
	ScreenManager.Render()
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
	ui.DrawAgentHUD(player.GetAgent())
	inventory := item.GetInventoryById(player.GetAgent().InventoryID)
	invCameraTransform := Cam.GetTransform().Inv()
	inventory.Draw(baseCtx, invCameraTransform)
	ui.DrawDamageIndicators(invCameraTransform)
	render.FlushUI()
}
