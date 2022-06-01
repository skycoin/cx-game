package spine

import (
	"math"
)

var (
	V0 = Vector{0, 0}
	V1 = Vector{1, 1}
)

type Vector struct {
	X, Y float32
}

func V(x, y float32) Vector { return Vector{x, y} }

func (a *Vector) Set(x, y float32) { a.X, a.Y = x, y }

func (a Vector) XY() (float32, float32) { return a.X, a.Y }

func (a Vector) V() struct{ X, Y float64 } {
	return struct{ X, Y float64 }{float64(a.X), float64(a.Y)}
}

func (a Vector) Add(b Vector) Vector  { return Vector{a.X + b.X, a.Y + b.Y} }
func (a Vector) Sub(b Vector) Vector  { return Vector{a.X - b.X, a.Y - b.Y} }
func (a Vector) EMul(b Vector) Vector { return Vector{a.X * b.X, a.Y * b.Y} }
func (a Vector) Len() float32         { return sqrt(a.Len2()) }
func (a Vector) Len2() float32        { return a.X*a.X + a.Y*a.Y }

type Transform struct {
	Translate Vector
	Rotate    float32
	Scale     Vector
	Shear     Vector
}

func NewTransform() Transform {
	transform := Transform{}
	transform.Translate = V0
	transform.Rotate = 0
	transform.Scale = V1
	transform.Shear = V0
	return transform
}

func (a Transform) Combine(b Transform) Transform {
	return Transform{
		Translate: a.Translate.Add(b.Translate),
		Rotate:    a.Rotate + b.Rotate,
		Scale:     a.Scale.EMul(b.Scale),
		Shear:     a.Shear.Add(b.Shear),
	}
}

type Affine struct {
	M00, M01, M02 float32
	M10, M11, M12 float32
}

func Identity() Affine {
	return Affine{
		1, 0, 0,
		0, 1, 0,
	}
}

func Translation(t Vector) Affine {
	return Affine{
		1, 0, t.X,
		0, 1, t.Y,
	}
}

func Rotation(r float32) Affine {
	sn, cs := sincos(r)
	return Affine{
		cs, -sn, 0,
		sn, cs, 0,
	}
}

func Scale(s Vector) Affine {
	return Affine{
		s.X, s.Y, 0,
		s.X, s.Y, 0,
	}
}

func (transform Transform) Affine() Affine {
	if transform.Shear == V0 {
		t, r, s := transform.Translate, transform.Rotate, transform.Scale
		switch r {
		case 0, 2 * math.Pi:
			return Affine{
				s.X * 1, -s.Y * 0, t.X,
				s.X * 0, s.Y * 1, t.Y,
			}
		case -math.Pi / 2:
			return Affine{
				s.X * 0, -s.Y * -1, t.X,
				s.X * -1, s.Y * 0, t.Y,
			}
		case math.Pi / 2:
			return Affine{
				s.X * 0, -s.Y * 1, t.X,
				s.X * 1, s.Y * 0, t.Y,
			}
		case -math.Pi, math.Pi:
			return Affine{
				s.X * -1, -s.Y * 0, t.X,
				s.X * 0, s.Y * -1, t.Y,
			}
		default:
			sn, cs := sincos(r)
			return Affine{
				s.X * cs, -s.Y * sn, t.X,
				s.X * sn, s.Y * cs, t.Y,
			}
		}
	} else {
		t, r, s, h := transform.Translate, transform.Rotate, transform.Scale, transform.Shear
		snx, csx := sincos(r + h.X)
		sny, csy := sincos(r + math.Pi/2 + h.Y)
		return Affine{
			s.X * csx, s.Y * sny, t.X,
			s.X * snx, s.Y * csy, t.Y,
		}
	}
}

func (a Affine) Mul(b Affine) Affine {
	m00 := a.M00*b.M00 + a.M01*b.M10
	m01 := a.M00*b.M01 + a.M01*b.M11
	m02 := a.M00*b.M02 + a.M01*b.M12 + a.M02
	m10 := a.M10*b.M00 + a.M11*b.M10
	m11 := a.M10*b.M01 + a.M11*b.M11
	m12 := a.M10*b.M02 + a.M11*b.M12 + a.M12

	return Affine{
		m00, m01, m02,
		m10, m11, m12,
	}
}

func (aff Affine) Transform(p Vector) Vector {
	return Vector{
		X: aff.M00*p.X + aff.M01*p.Y + aff.M02*1,
		Y: aff.M10*p.X + aff.M11*p.Y + aff.M12*1,
	}
}

func (aff Affine) Translation() Vector {
	return Vector{
		X: aff.M02,
		Y: aff.M12,
	}
}

func (aff Affine) WeightedTransform(w float32, p Vector) Vector {
	return Vector{
		X: w * (aff.M00*p.X + aff.M01*p.Y + aff.M02*1),
		Y: w * (aff.M10*p.X + aff.M11*p.Y + aff.M12*1),
	}
}

// Conversion from an to different formats
func (a Affine) Row32() [6]float32 {
	return [6]float32{a.M00, a.M01, a.M02, a.M10, a.M11, a.M12}
}

func (a Affine) Row64() [6]float64 {
	return [6]float64{float64(a.M00), float64(a.M01), float64(a.M02), float64(a.M10), float64(a.M11), float64(a.M12)}
}

func (a Affine) Col32() [6]float32 {
	return [6]float32{a.M00, a.M10, a.M01, a.M11, a.M02, a.M12}
}

func (a Affine) Col64() [6]float64 {
	return [6]float64{float64(a.M00), float64(a.M10), float64(a.M01), float64(a.M11), float64(a.M02), float64(a.M12)}
}

// Conversion from an to different formats
func (a Affine) Row32s() (m00, m01, m02, m10, m11, m12 float32) {
	return a.M00, a.M01, a.M02, a.M10, a.M11, a.M12
}

func (a Affine) Row64s() (m00, m01, m02, m10, m11, m12 float64) {
	return float64(a.M00), float64(a.M01), float64(a.M02), float64(a.M10), float64(a.M11), float64(a.M12)
}

func (a Affine) Col32s() (m00, m10, m01, m11, m02, m12 float32) {
	return a.M00, a.M10, a.M01, a.M11, a.M02, a.M12
}

func (a Affine) Col64s() (m00, m10, m01, m11, m02, m12 float64) {
	return float64(a.M00), float64(a.M10), float64(a.M01), float64(a.M11), float64(a.M02), float64(a.M12)
}
