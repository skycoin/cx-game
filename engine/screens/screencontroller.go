package screens

import (
	"github.com/skycoin/cx-game/engine/input/inputhandler"
)

type SCREEN int

const (
	TITLE SCREEN = iota
	GAME
	MENU
	MAP
	/*
		menu
		options screen
		player/norma/game screen
		map screen
	*/
)

type ScreenManager struct {
	screenHandlers map[SCREEN]ScreenHandler
	//todo have a queue of active screens
	activeScreenIndex SCREEN
}

func NewDevScreenManager() *ScreenManager {
	newScreenController := ScreenManager{
		screenHandlers: make(map[SCREEN]ScreenHandler),
	}
	// newScreenController.screenHandlers[TITLE] = NewTitleScreenHandler()
	newScreenController.screenHandlers[GAME] = NewGameScreenHandler()
	newScreenController.activeScreenIndex = GAME
	// newScreenController.screenHandlers[MENU] = NewMenuScreenHandler()
	// newScreenController.screenHandlers[MAP] = NewMapScreenHandler()
	return &newScreenController
}

// func (s *ScreenManager) ProcessInput() {
// 	//todo multiple active screens
// 	s.screenHandlers[s.activeScreenIndex].ProcessInput()

// 	//reset event queue
// 	input.InputEvents = input.InputEvents[:0]
// }

func (s *ScreenManager) Update(dt float32) {
	s.screenHandlers[s.activeScreenIndex].Update(dt)
}

func (s *ScreenManager) FixedUpdate() {
	s.screenHandlers[s.activeScreenIndex].FixedUpdate()
}
func (s *ScreenManager) Render(ctx DrawContext) {
	s.screenHandlers[s.activeScreenIndex].Render(ctx)
}

func (s *ScreenManager) RegisterButton(screen SCREEN, info inputhandler.ActionInfo) {
	s.screenHandlers[screen].RegisterAction(info)
}
