package screens

import "github.com/skycoin/cx-game/engine/input/inputhandler"

type ScreenHandler interface {
	Update(float32)
	FixedUpdate()
	Render(DrawContext)
	RegisterAction(inputhandler.ActionInfo)
}
