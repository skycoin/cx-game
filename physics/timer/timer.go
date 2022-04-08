package timer

import (
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath/math32"
)

var Accumulator float32

func GetTimeBetweenTicks() float32 {
	return math32.Mod(Accumulator, constants.MS_PER_TICK)
}
