package main

import (
	"time"

	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/sound"
)

func main() {
	sound.Init()
	sound.LoadSound("boo", "boo.wav")
	enemyPos := physics.Vec2{-6, 2}
	sound.Play2DSound("boo", &enemyPos, sound.SoundOptions{
		IsStatic: true,
	})
	//imitate game loop
	pos := physics.Vec2{3, 0}
	sound.SetListenerPosition(pos)
	for {
		enemyPos.X += 0.3
		time.Sleep(50 * time.Millisecond)
	}
}
