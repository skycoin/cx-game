package input

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/cxmath"
)

var (
	MouseCoords                 cxmath.Vec2
	MousePressed                bool = false
	physicalToViewPortTransform mgl32.Mat4
)

func SetPhysicalToViewportTransform(transform mgl32.Mat4) {
	physicalToViewPortTransform = transform
}
