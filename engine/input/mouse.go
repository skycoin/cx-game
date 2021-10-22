package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/cxmath"
)

var (
	physicalToViewPortTransform mgl32.Mat4

	//list of mouse buttons pressed
	//mousebutton1 = LMB
	//mousebutton2 = RMB
	mousePressed map[glfw.MouseButton]bool = make(map[glfw.MouseButton]bool)

	//for physics
	FixedMouseInfo MouseInfo
	RenderMouseInfo      MouseInfo
)

type MouseInfo struct {
	PrevPos cxmath.Vec2
	Pos     cxmath.Vec2
	Offset  cxmath.Vec2
}

func (m *MouseInfo) SetPos(x, y float64) {
	m.PrevPos = m.Pos
	m.Pos.X, m.Pos.Y = float32(x), float32(y)

	m.Offset = m.Pos.Sub(m.PrevPos)
}

func SetPhysicalToViewportTransform(transform mgl32.Mat4) {
	physicalToViewPortTransform = transform

}

func PollMouse(window *glfw.Window, fixed bool) {
	x, y := window.GetCursorPos()

	if fixed {
		FixedMouseInfo.SetPos(x, y)
	} else {
		RenderMouseInfo.SetPos(x, y)
	}
}

func mouseButtonCallback(
	w *glfw.Window, b glfw.MouseButton, a glfw.Action, mk glfw.ModifierKey,
) {
	if a == glfw.Press {
		mousePressed[b] = true
	} else if a == glfw.Release {
		mousePressed[b] = false
	}
}
