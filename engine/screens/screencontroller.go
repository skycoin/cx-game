package screens

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
	activeScreen SCREEN
}

func NewDevScreenManager() *ScreenManager {
	newScreenController := ScreenManager{
		screenHandlers: make(map[SCREEN]ScreenHandler),
	}

	// newScreenController.screenHandlers[TITLE] = NewTitleScreenHandler()

	newScreenController.screenHandlers[GAME] = NewGameScreenHandler()
	newScreenController.activeScreen = GAME
	// newScreenController.screenHandlers[MENU] = NewMenuScreenHandler()
	// newScreenController.screenHandlers[MAP] = NewMapScreenHandler()
	return &newScreenController
}
 
