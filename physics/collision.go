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

func CheckCollisions(from,to mgl32.Vec2) (collision Collision, collided bool) {
	at := from
	distSqr := to.Sub(from).LenSqr()
	dir := to.Sub(from).Normalize()
	for at.Sub(from).LenSqr() < distSqr {
		collision,collided :=
			CheckCollision(mgl32.Translate3D(at.X(),at.Y(),0))
		if collided { return collision,collided }
		at = at.Add(dir)
	}
	return Collision{}, false
}

func CheckCollision(transform mgl32.Mat4) (collision Collision, collided bool) {
	for _,body := range bodies {
		if body.CollidesWith(transform) {
			return Collision{Body:body}, true
		}
	}
	return Collision{}, false
}
