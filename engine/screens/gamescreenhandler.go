package screens

import "github.com/skycoin/cx-game/engine/input/inputhandler"

type GameScreenHandler struct {
	inputController inputhandler.InputController
}

func NewGameScreenHandler() *GameScreenHandler {
	newScreenHandler := GameScreenHandler{
		inputHandler: inputhandler.NewAgentInputController(),
		
	}

	return &newScreenHandler
}

func (g *GameScreenHandler) Render() {

}

func (g *GameScreenHandler) SetInputHandler() {

}
