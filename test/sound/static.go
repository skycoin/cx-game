package main

import (
	"time"

	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/sound"
)

func main() {
	sound.Init()

	sound.LoadSound("boo", "jump.wav")

	sound.PlaySound("boo")

	pos := physics.Vec2{-5, 0}
	for {
		sound.Update()
		pos.X += 0.01
		sound.SetListenerPosition(pos)
		time.Sleep(50 * time.Millisecond)
	}
}
