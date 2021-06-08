package main

import (
	"time"

	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/sound"
)

func main() {
	sound.Init()

	sound.LoadSound("boo", "boo.wav")

	enemyPos := physics.Vec2{-3, 2}
	// sound.Play3DSound("bloop", physics.Vec2{-3, 2})
	// sound.PlaySound("piano")
	sound.Play2DSound("boo", &enemyPos, true)

	//imitate game loop
	pos := physics.Vec2{-6, 0}

	for {
		sound.SetListenerPosition(pos)
		pos.X += 0.3
		time.Sleep(50 * time.Millisecond)
	}
}
