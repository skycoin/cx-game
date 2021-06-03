package physics

import "math"

type Vec2 struct {
	X, Y float32
}

type Vec2i struct {
	X, Y int32
}

func (v1 Vec2) Add(v2 Vec2) Vec2 {
	v1.X += v2.X
	v1.Y += v2.Y

	return v1
}

func (v1 Vec2) Sub(v2 Vec2) Vec2 {
	return v1.Add(v2.Mult(-1))
}

func (v1 Vec2) Mult(n float32) Vec2 {
	v1.X *= n
	v1.Y *= n

	return v1
}

func (v1 Vec2) LengthSqr() float32 {
	return v1.X*v1.X + v1.Y*v1.Y
}

func (v1 Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v1.LengthSqr())))
}

func (v1 Vec2) IsZero() bool {
	return v1.X == 0 && v1.Y == 0
}

func (v1 Vec2) Normalize() Vec2 {
	if !v1.IsZero() {
		length := v1.Length()
		v1.X = v1.X / length
		v1.Y = v1.Y / length
	}

	return v1
}

func (v1 Vec2i) Mult(n int32) Vec2i{
	return Vec2i {
		X: v1.X*n,
		Y: v1.Y*n,
	}
}

func (v1 Vec2i) Add(v2 Vec2i) Vec2i {
	return Vec2i {
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}
