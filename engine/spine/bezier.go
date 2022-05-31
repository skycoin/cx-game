package spine

import (
	"math"
)

var (
	linearTag = math.Float32frombits(0x7f800000) // +Inf
	stepTag   = math.Float32frombits(0xff800000) // -Inf
)

type Curve struct {
	Linear bool
	Step   bool
	Bezier Bezier01D
}

// Bezier01 calculates point where P0 = {0,0} and P3 = {1,1}
func Bezier01(p1, p2 Vector, t float32) Vector {
	// B(t) = (1-t)^3*P0 + 3(1-t)^2*t*P1 + 3(1-t)*t^2*P2 + t^3*P3
	//   because P0 == {0,0} and P3 == {1,1}
	// B(t) = 3(1-t)^2*t*P1 + 3(1-t)*t^2*P2 + V(t^3)
	return Vector{
		X: 3*(1-t)*(1-t)*t*p1.X + 3*(1-t)*t*t*p2.X + t*t*t,
		Y: 3*(1-t)*(1-t)*t*p1.Y + 3*(1-t)*t*t*p2.Y + t*t*t,
	}
}

func (curve *Curve) setTag(tag float32) { curve.Bezier.D.X = tag }
func (curve *Curve) getTag() float32    { return curve.Bezier.D.X }

func (curve *Curve) SetLinear()              { curve.setTag(linearTag) }
func (curve *Curve) SetStep()                { curve.setTag(stepTag) }
func (curve *Curve) SetBezier(p1, p2 Vector) { curve.Bezier.Set(p1, p2) }

func (curve *Curve) Evaluate(p float32) float32 {
	tag := curve.getTag()
	if tag == linearTag {
		return p
	} else if tag == stepTag {
		return 0
	}
	return curve.Bezier.Y(p)
}

// max deviation ~0.1, ~30cy/Y
type Bezier01D struct {
	D   Vector
	DD  Vector
	DDD Vector
}

func NewBezier01D(p1, p2 Vector) Bezier01D {
	var bez Bezier01D
	bez.Set(p1, p2)
	return bez
}

func (bez *Bezier01D) Set(p1, p2 Vector) {
	const S = 1.0 / 10.0
	var t0, t1 Vector

	t0.X = -p1.X*2 + p2.X
	t0.Y = -p1.Y*2 + p2.Y
	t1.X = (p1.X-p2.X)*3 + 1
	t1.Y = (p1.Y-p2.Y)*3 + 1

	bez.D.X = p1.X*3*S + t0.X*3*S*S + t1.X*S*S*S
	bez.D.Y = p1.Y*3*S + t0.Y*3*S*S + t1.Y*S*S*S
	bez.DD.X = t0.X*6*S*S + t1.X*6*S*S*S
	bez.DD.Y = t0.Y*6*S*S + t1.Y*6*S*S*S
	bez.DDD.X = t1.X * 6 * S * S * S
	bez.DDD.Y = t1.Y * 6 * S * S * S
}

func (bez *Bezier01D) Y(x float32) float32 {
	x = clamp01(x)

	d, dd, ddd := bez.D, bez.DD, bez.DDD
	b := d
	for i := 0; i < 8; i++ {
		if b.X >= x {
			last := b.Sub(d)
			return last.Y + (b.Y-last.Y)*(x-last.X)/(b.X-last.X)
		}
		d = d.Add(dd)
		dd = dd.Add(ddd)
		b = b.Add(d)
	}
	if b.X >= x {
		last := b.Sub(d)
		return last.Y + (b.Y-last.Y)*(x-last.X)/(b.X-last.X)
	}
	return b.Y + (1-b.Y)*(x-b.X)/(1-b.X)
}

// max deviation ~0.3, ~4.7cy/Y
type Bezier01T struct {
	A, B Vector
}

func NewBezier01T(p1, p2 Vector) Bezier01T {
	var bez Bezier01T
	bez.Set(p1, p2)
	return bez
}

func (bez *Bezier01T) Set(p1, p2 Vector) {
	// B(t) = (1-t)^3*P0 + 3(1-t)^2*t*P1 + 3(1-t)*t^2*P2 + t^3*P3
	// B(t) = 3(1-t)^2*t*P1 + 3(1-t)*t^2*P2 + V(t^3)
	const A = 1.0 / 3.0
	const B = 2.0 / 3.0

	bez.A.X = 3*(1-A)*(1-A)*A*p1.X + 3*(1-A)*A*A*p2.X + A*A*A
	bez.A.Y = 3*(1-A)*(1-A)*A*p1.Y + 3*(1-A)*A*A*p2.Y + A*A*A

	bez.B.X = 3*(1-B)*(1-B)*B*p1.X + 3*(1-B)*B*B*p2.X + B*B*B
	bez.B.Y = 3*(1-B)*(1-B)*B*p1.Y + 3*(1-B)*B*B*p2.Y + B*B*B
}

func (bez *Bezier01T) Y(x float32) float32 {
	a, b := bez.A, bez.B
	if x < a.X {
		mix := (a.X - x) / a.X
		return lerp(0, a.Y, 1-mix)
	} else if x < b.X {
		mix := (b.X - x) / (b.X - a.X)
		return lerp(a.Y, b.Y, 1-mix)
	} else if x < 1 {
		mix := (1 - x) / (1 - b.X)
		return lerp(b.Y, 1, 1-mix)
	}
	return 1
}
