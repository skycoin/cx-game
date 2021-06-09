package main

import (
	"fmt"
	"time"

	"github.com/ikemen-engine/go-openal/openal"
	"github.com/skycoin/cx-game/sound"
)

func main() {
	sound.Init()

	sound.LoadSound("boo", "boo.wav")

	// sound.Play3DSound("boo", &physics.Vec2{-3, 2}, true)
	// sound.PlaySound("piano")
	buffer, _ := sound.NewBuffer("boo.wav")
	source2 := openal.NewSource()
	source := openal.NewSource()

	listener := openal.Listener{}

	var vec openal.Vector
	listener.GetPosition(&vec)

	listener.SetGain(0)
	fmt.Println(vec)

	source.SetBuffer(*buffer)
	source2.SetBuffer(*buffer)

	source.Play()
	time.Sleep(100 * time.Millisecond)

	source2.Play()

	// //imitate game loop
	// pos := physics.Vec2{-6, 0}
	// var posRetrieved openal.Vector
	volume := float32(0)
	for {
		listener.SetGain(volume)
		volume += 0.03
		fmt.Println(volume)
		// sound.SetListenerPosition(pos)
		// sound.Listener.GetPosition(&posRetrieved)
		// fmt.Println(posRetrieved)
		// pos.X += 0.5
		time.Sleep(50 * time.Millisecond)
	}
}
