package inputhandler

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type InputMapper struct {
	enumToKey map[int]glfw.Key
	keyToEnum map[glfw.Key]int
}

func (i InputMapper) GetKey(enum int) glfw.Key {
	key, ok := i.enumToKey[enum]

	if !ok {
		log.Fatalln("No such enum!")
	}

	return key
}

func (i InputMapper) BindKey(key glfw.Key, enum int) {
	if _, ok := i.enumToKey[enum]; ok {
		//key is already bind, add logging
		log.Println("REBINDING KEY")
	}
	i.enumToKey[enum] = key

	//check if key is bind,
	if _, ok := i.keyToEnum[key]; ok {
		oldEnum := i.keyToEnum[key]

		//reset old keybind
		//todo add proper logging here
		log.Println("RESETTING OLD KEYBIND")
		i.enumToKey[oldEnum] = glfw.KeyUnknown
	}

	i.keyToEnum[key] = enum
}
