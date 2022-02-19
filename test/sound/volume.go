package main

import (
	"fmt"
	"time"

	"github.com/skycoin/cx-game/lib/openal"
	"github.com/skycoin/cx-game/engine/sound"
)

func main() {
	sound.Init()
	sound.LoadSound("boo", "boo.wav")
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
	volume := float32(0)
	//imitate game loop
	for {
		listener.SetGain(volume)
		volume += 0.03
		fmt.Println(volume)
		time.Sleep(50 * time.Millisecond)
	}
}
