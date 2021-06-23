package worldcollider

import (
	"github.com/go-gl/mathgl/mgl32"
)

type WorldCollider interface {
	TileIsSolid(x,y int) bool
	WrapAround(rawPosition mgl32.Vec2) (wrappedPosition mgl32.Vec2)
}
