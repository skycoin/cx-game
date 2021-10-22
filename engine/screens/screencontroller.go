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

type ScreenController struct {
	screenHandlers map[SCREEN]ScreenHandler
}

func NewDevScreenController() *ScreenController {
	newScreenController := ScreenController{
		screenHandlers: make(map[SCREEN]ScreenHandler),
	}

	newScreenController.screenHandlers[TITLE] = NewTitleScreenHandler()

	newScreenController.screenHandlers[GAME] = NewGameScreenHandler()
	newScreenController.screenHandlers[MENU] = NewMenuScreenHandler()
	newScreenController.screenHandlers[MAP] = NewMapScreenHandler()

}
