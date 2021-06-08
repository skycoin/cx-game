package cxmath

type Vec2i struct {
	X, Y int32
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
