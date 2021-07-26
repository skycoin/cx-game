package cxmath

import (
	"github.com/skycoin/cx-game/cxmath/mathi"
)

type Frustum struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

func (this Frustum) Intersect (other Frustum) Frustum {
	return Frustum {
		Left: mathi.Max(this.Left, other.Left),
		Right: mathi.Min(this.Right, other.Right),
		Top: mathi.Min(this.Top, other.Top),
		Bottom: mathi.Max(this.Bottom, other.Bottom),
	}
}
