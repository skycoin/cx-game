package verlet

import (
	"github.com/go-gl/mathgl/mgl32"
)

// TODO get some bounces in here

// Newtonian physics with verlet integration
// https://en.wikipedia.org/wiki/Verlet_integration

// 2D verlet state
type Verlet2 struct { Position,Velocity mgl32.Vec2 }

func NewVerlet2(position, velocity mgl32.Vec2) Verlet2 {
	return Verlet2 { Position: position, Velocity: velocity }
}

func (verlet *Verlet2) Integrate(dt float32, a mgl32.Vec2) {
	// copy old state
	p := verlet.Position
	v := verlet.Velocity
	// overwrite verlet state
	verlet.Position = p.Add(v.Mul(dt)).Add(a.Mul(0.5*dt*dt))
	verlet.Velocity = v.Add(a.Mul(0.5*dt))
}

