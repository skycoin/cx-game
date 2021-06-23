package physics

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Collision struct {
	Body *Body
}

func (body *Body) CollidesWith(other mgl32.Mat4) bool {
	// TODO something better
	disp := other.Col(3).Vec2().Sub(body.Transform().Col(3).Vec2())
	return disp.LenSqr() <= 0.5*0.5

}

func CheckCollision(transform mgl32.Mat4) (collision Collision, collided bool) {
	for _,body := range bodies {
		if body.CollidesWith(transform) {
			return Collision{Body:body}, true
		}
	}
	return Collision{}, false
}
