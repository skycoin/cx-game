package game

import (
	"github.com/skycoin/cx-game/engine/screens"
)

var debugTileInfo bool = true

func Draw() {
	ScreenManager.Render(screens.DrawContext{
		Cam:     Cam,
		World:   &World,
		Window:  &win,
		Console: Console,
		Fps:     fps,
		Player:  player,
	})
}

// func DrawUI() {
// 	baseCtx := win.DefaultRenderContext()
// 	topLeftCtx :=
// 		render.CenterToTopLeft(baseCtx)
// 	if debugTileInfo {
// 		ui.DrawString(
// 			tileText,
// 			mgl32.Vec4{0.3, 0.9, 0.4, 1},
// 			ui.AlignCenter,
// 			topLeftCtx.PushLocal(mgl32.Translate3D(25, -5, 0)),
// 		)
// 	}
// 	ui.DrawAgentHUD(player.GetAgent())
// 	inventory := item.GetInventoryById(player.GetAgent().InventoryID)
// 	invCameraTransform := Cam.GetTransform().Inv()
// 	inventory.Draw(baseCtx, invCameraTransform)
// 	ui.DrawDamageIndicators(invCameraTransform)
// 	render.FlushUI()
// }
