package console

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/ui"
	"github.com/skycoin/cx-game/render"
)

type Console struct {
	active bool
	line   string
	output string
}

func New() Console {
	return Console{}
}

func (console *Console) IsActive() bool { return console.active }

func (console *Console) OnChar(w *glfw.Window, char rune) {
	if string(char) == "\n" || string(char) == "`" {
		return
	}
	console.line = console.line + string(char)
}

func (console *Console) ToggleActive(window *glfw.Window) {
	console.active = !console.active
	if console.active {
		window.SetCharCallback(console.OnChar)
		console.line = ""
	} else {
		window.SetCharCallback(nil)
	}
}

func (console *Console) Update(window *glfw.Window, ctx CommandContext) {
	if input.GetKeyDown(glfw.KeyGraveAccent) {
		console.ToggleActive(window)
	}
	if console.active {
		if input.GetKeyDown(glfw.KeyBackspace) && len(console.line) > 0 {
			console.line = console.line[:len(console.line)-1]
		}
		if input.GetKeyDown(glfw.KeyEnter) {
			console.Command(ctx)
		}
	}
}

func (console *Console) Draw(ctx render.Context) {
	if !console.active {
		return
	}
	ctx = render.CenterToTopLeft(ctx).
		PushLocal(mgl32.Translate3D(1, -10, 0))
	ui.DrawString(
		"> "+console.line, mgl32.Vec4{1, 0, 0, 1},
		ui.AlignLeft,
		ctx,
	)
	ui.DrawString(
		console.output, mgl32.Vec4{0, 1, 0, 1},
		ui.AlignLeft,
		ctx.PushLocal(mgl32.Translate3D(0,1,0)),
	)
}

func (console *Console) Command(ctx CommandContext) {
	console.output = processCommand(console.line, ctx)
	console.line = ""
}
