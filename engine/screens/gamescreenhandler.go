package screens

import "github.com/skycoin/cx-game/engine/input/inputhandler"

type GameScreenHandler struct {
	inputHandler *inputhandler.AgentInputHandler
}

func NewGameScreenHandler() *GameScreenHandler {
	newScreenHandler := GameScreenHandler{
		inputHandler: inputhandler.NewAgentInputHandler(),
	}

	return &newScreenHandler
}

func (g *GameScreenHandler) ProcessInput() {
	// g.inputHandler.ProcessEvents()
	g.inputHandler.UpdateKeyState()
}

func (g *GameScreenHandler) Update() {

}

func (g *GameScreenHandler) Render() {

}
