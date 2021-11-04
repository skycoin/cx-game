package screens

import "github.com/skycoin/cx-game/engine/input/inputhandler"

type ScreenHandler interface {
	ProcessInput()
	Update(float32)
	FixedUpdate()
	Render()
	RegisterAction(inputhandler.ActionInfo)
}
