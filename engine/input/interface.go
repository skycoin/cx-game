package input

import (
	"fmt"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
)

var DEBUG = false

//continuos keys, holding
func GetButton(button string) bool {
	keyComb, ok := ActiveButtonsToKeys[button]
	if !ok && DEBUG {
		log.Printf("BUTTON [%s] IS NOT MAPPED!\n", button)
		return false
	}
	pressed, ok := keyPressed[keyComb.key]
	if !ok {
		// log.Printf("ERROR!")
		return false
	}
	if modifierKey == keyComb.modifierKey {
		return pressed
	}
	return false
}

//action keys, if pressed once
func GetButtonDown(button string) bool {
	keyComb, ok := ActiveButtonsToKeys[button]
	if !ok && DEBUG {
		log.Printf("BUTTON [%s] IS NOT MAPPED\n", button)
		return false
	}
	pressed, ok := keyPressedDown[keyComb.key]
	if !ok {
		return false
	}
	if modifierKey == keyComb.modifierKey && pressed {
		keyPressedDown[keyComb.key] = false
		return true

	}

	return false
}

func GetButtonUp(button string) bool {
	keyComb, ok := ActiveButtonsToKeys[button]
	if !ok && DEBUG {
		log.Printf("BUTTON [%s] IS NOT MAPPED\n", button)
		return false
	}
	if modifierKey == keyComb.modifierKey {

		return GetKeyDown(keyComb.key)
	}
	return false
}

func GetKey(key glfw.Key) bool {
	if modifierKey == 0 {
		return keyPressed[key]
	}
	return false
}

func GetKeyDown(key glfw.Key) bool {
	pressed, ok := keysPressedUp[key]
	if !ok {
		return false
	}
	keysPressedUp[key] = false
	return pressed
}

func GetKeyIsUp(key glfw.Key) bool {
	return window_.GetKey(key) == glfw.Press
}

func GetLastKey() glfw.Key {
	key := lastKeyPressed
	//consume and set to false
	lastKeyPressed = glfw.KeyUnknown
	return key
}

func GetAxis(axis Axis) float32 {
	if axis == HORIZONTAL {
		return cxmath.BoolToFloat(GetButton("right")) - cxmath.BoolToFloat(GetButton("left"))
	} else { // VERTICAL
		return cxmath.BoolToFloat(GetButton("up")) - cxmath.BoolToFloat(GetButton("down"))
	}

}

func GetMouseX() float32 {
	return float32(MouseCoords.X)
}
func GetMouseY() float32 {
	return float32(MouseCoords.Y)
}

func GetMousePos() mgl32.Vec2 {
	physicalX := float32(MouseCoords.X)
	physicalY := float32(MouseCoords.Y)

	physicalPos := mgl32.Vec2{physicalX, physicalY}
	physicalPosHomogenous :=
		mgl32.Vec4{physicalPos.X(), physicalPos.Y(), 0, 1}

	transform := physicalToViewPortTransform
	virtualPos := transform.Mul4x1(physicalPosHomogenous).Vec2()

	return virtualPos
}

func Reset() {
	//reset lastkeyPressed
	return
}

//for debug purposes

type Flag struct {
	enableFunc  func()
	disableFunc func()
	description string
	isEnabled   bool
}

type Slider struct {
	value       *float32
	step        float32
	min         float32
	max         float32
	incKey      glfw.Key
	decKey      glfw.Key
	mKey        glfw.ModifierKey
	description string
}

var registeredFlags = make(map[glfw.ModifierKey]map[glfw.Key]*Flag)

var registeredSliders = make([]Slider, 0)

func RegisterFlag(enableFunc, disableFunc func(), switchKey glfw.Key, mk glfw.ModifierKey, description string) {
	if registeredFlags[mk] == nil {
		registeredFlags[mk] = make(map[glfw.Key]*Flag)
	}
	_, ok := registeredFlags[mk][switchKey]
	if !ok {
		registeredFlags[mk][switchKey] = &Flag{
			enableFunc:  enableFunc,
			disableFunc: disableFunc,
			isEnabled:   false,
			description: description,
		}
		return
	}
	log.Fatalf("Key is already reserved")
}

func ProcessFlags(key glfw.Key, mk glfw.ModifierKey) {

	flagg, ok := registeredFlags[mk][key]
	if ok {
		flagg.isEnabled = !flagg.isEnabled
		fmt.Printf("%v: %v\n", flagg.description, flagg.isEnabled)

		if flagg.isEnabled {
			flagg.enableFunc()
		} else {
			flagg.disableFunc()
		}
	}
}

func ProcessSliders(key glfw.Key, mKey glfw.ModifierKey) {
	for _, slider := range registeredSliders {
		if slider.mKey == 0 || slider.mKey == mKey {
			if key == slider.decKey || key == slider.incKey {
				inc := slider.step
				if key == slider.decKey {
					inc *= -1
				}
				*slider.value = math32.Clamp(*slider.value+inc, slider.min, slider.max)
				fmt.Printf("%v:  %v\n", slider.description, *slider.value)
				return
			}
		}
	}
}

func RegisterSlider(
	value *float32, step, min, max float32, incKey, decKey glfw.Key,
	mKey glfw.ModifierKey, description string) {
	newSlider := Slider{
		value:       value,
		step:        step,
		min:         min,
		max:         max,
		incKey:      incKey,
		decKey:      decKey,
		mKey:        mKey,
		description: description,
	}

	registeredSliders = append(registeredSliders, newSlider)
}
