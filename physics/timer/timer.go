package timer

import (
	"github.com/skycoin/cx-game/constants/physicsconstants"
	"github.com/skycoin/cx-game/cxmath/math32"
)

var Accumulator float32

func GetTimeBetweenTicks() float32 {
	return math32.Mod(Accumulator, physicsconstants.PHYSICS_TIMESTEP)
}
