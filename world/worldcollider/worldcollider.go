package worldcollider

import (
	"github.com/go-gl/mathgl/mgl32"
)

type WorldCollider interface {
	TileIsSolid(x,y int) bool
	TileTopIsSolid(x,y int, ignorePlatforms bool) bool
	WrapAroundOffset(rawPosition mgl32.Vec2) (wrappedPosition mgl32.Vec2)
}
