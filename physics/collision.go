package physics

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Collision struct {}

func CheckCollision(transform mgl32.Mat4) (collision Collision, collided bool) {
	return Collision{}, false
}
