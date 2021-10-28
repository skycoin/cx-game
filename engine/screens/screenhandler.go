package screens

type ScreenHandler interface {
	ProcessInput()
	Update()
	Render()
}
