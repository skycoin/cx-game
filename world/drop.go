// dictates which items drop from a tile
package world

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/types"
)

type Drop struct {
	Item        types.ItemTypeID
	Count       int
	Probability float32
}

type Drops []Drop

// generate a drop by sampling the drop entries in this list
func (drops Drops) Drop() Drop {
	// for this to work properly,
	// it is assumed that the sum of the drop probabilities is < 1
	p := rand.Float32()
	for _, drop := range drops {
		p -= drop.Probability
		if p < 0 {
			return drop
		}
	}
	return Drop{} // drop nothing
}
