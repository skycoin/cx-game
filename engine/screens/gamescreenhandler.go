package screens

import "github.com/skycoin/cx-game/engine/input/inputhandler"

type GameScreenHandler struct {
	inputHandler inputhandler.InputHandler
}

func NewGameScreenHandler() *GameScreenHandler {
	newScreenHandler := GameScreenHandler{}

	return &newScreenHandler
}

func (g *GameScreenHandler) Render() {

}

func (g *GameScreenHandler) SetInputHandler() {

}
