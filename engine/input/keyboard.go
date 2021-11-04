package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	//for actions
	KeyPressed = make(map[glfw.Key]bool)

	modifierKey glfw.ModifierKey

	//for each context, have map of registered buttons to keys
	ButtonsMap          = make(map[InputContext]map[string]KeyComb)
	ActiveButtonsToKeys map[string]KeyComb
	lastKeyPressed      glfw.Key

	InputEvents []EventInfo
)

type EventInfo struct {
	Type   EventType
	Button int
}

type EventType int

const (
	KEY_DOWN EventType = iota
	KEY_UP

	//mouse events
	MOUSE_DOWN
	MOUSE_UP

	MOUSE_SCROLL_UP
	MOUSE_SCROLL_DOWN

	MOUSE_MOVED
)

type Axis int

type KeyComb struct {
	key         glfw.Key
	modifierKey glfw.ModifierKey
}

const (
	HORIZONTAL Axis = iota
	VERTICAL
)

func keyCallback(
	w *glfw.Window,
	key glfw.Key, scancode int, action glfw.Action, mk glfw.ModifierKey,
) {
	modifierKey = mk

	newEvent := EventInfo{
		Button: int(key),
	}
	switch action {
	case glfw.Press:
		if key == glfw.KeyEscape {
			w.SetShouldClose(true)
		}
		KeyPressed[key] = true
		newEvent.Type = KEY_DOWN
		lastKeyPressed = key
		//todo check for maybe some concurrency issues, since filling and emptying array happens in different threads

		InputEvents = append(InputEvents, newEvent)
	case glfw.Release:
		KeyPressed[key] = false
		newEvent.Type = KEY_UP
		InputEvents = append(InputEvents, newEvent)
	}

}
