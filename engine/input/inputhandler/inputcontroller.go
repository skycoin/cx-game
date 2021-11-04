package inputhandler

type InputController struct {
	inputHandlers      []InputHandler
	activeInputHandler InputHandler
}

func NewInputController() *InputController {
	newInputController := InputController{
		inputHandlers: make([]InputHandler, 0),
	}

	newInputController.inputHandlers = append(
		newInputController.inputHandlers,
		NewAgentInputHandler(),
	)
	newInputController.activeInputHandler = NewAgentInputHandler()
	// newInputController.SetActiveInputHandler()

	return &newInputController
}
