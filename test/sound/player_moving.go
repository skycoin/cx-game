package main

import (
	"time"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/sound"
)

func main() {
	sound.Init()
	sound.LoadSound("boo", "boo.wav")
	enemyPos := cxmath.Vec2{-3, 2}
	sound.Play2DSound("boo", &enemyPos, sound.SoundOptions{
		IsStatic: true,
		Gain:     3.5,
		Pitch:    1.5,
	})

	pos := cxmath.Vec2{-6, 0}

	for {
		sound.SetListenerPosition(pos)
		pos.X += 0.3
		time.Sleep(50 * time.Millisecond)
	}
}
