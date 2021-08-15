package particle_emitter

import "github.com/skycoin/cx-game/world"

func Update(dt float32, currentWorld *world.World) {
	for i := range oxygenEmitters {
		oxygenEmitters[i].Update(dt, currentWorld)
	}
}
