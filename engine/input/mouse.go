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
	FixedMouseInfo  MouseInfo
	RenderMouseInfo MouseInfo
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
	//todo add support for modifier keys
	if a == glfw.Press {
		mousePressed[b] = true
		InputEvents = append(InputEvents, EventInfo{
			Type:   MOUSE_DOWN,
			Button: int(b),
		})
	} else if a == glfw.Release {
		mousePressed[b] = false
		InputEvents = append(InputEvents, EventInfo{
			Type:   MOUSE_UP,
			Button: int(b),
		})
	}
}

func cursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	//convenient for when we need to turn off cursor after some idle period
	InputEvents = append(InputEvents, EventInfo{
		Type: MOUSE_MOVED,
	})
}

func scrollCallback(w *glfw.Window, xOff, yOff float64) {
	var scrollEvent EventType
	if yOff > 1 {
		scrollEvent = MOUSE_SCROLL_DOWN
	} else {
		scrollEvent = MOUSE_SCROLL_UP
	}

	InputEvents = append(InputEvents, EventInfo{
		Type: scrollEvent,
	})
}

func sizeCallback(w *glfw.Window, width, height int) {
	//todo
}
