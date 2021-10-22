package input

import (
	"fmt"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/skycoin/cx-game/cxmath/math32"
)

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
