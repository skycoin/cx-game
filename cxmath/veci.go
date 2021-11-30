package cxmath

import (
	"math"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/cxmath/mathi"
)

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

func (v1 Vec2i) Sub(v2 Vec2i) Vec2i {
	return v1.Add(v2.Mult(-1))
}

func (v1 Vec2i) Length() float32 {
	return float32(math.Sqrt(float64(v1.X + v1.Y)))
}

func (v1 Vec2i) ManhattanDist() int32 {
	return int32(mathi.Abs(int(v1.X)) + mathi.Abs(int(v1.Y)))
}

func (v1 Vec2i) Vec2() mgl32.Vec2 {
	return mgl32.Vec2 { float32(v1.X), float32(v1.Y) }
}

func (v1 Vec2i) Eq(v2 Vec2i) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}
