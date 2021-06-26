package math32

import (
	"math"
)

func Sign(x float32) float32 {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func Min(x,y float32) float32 {
	if x < y {
		return x
	} else {
		return y
	}
}

func Abs(x float32) float32 {
	if x < 0 { x *= -1 }
	return x
}

// return the value with the lesser magnitude.
// result maintains original sign.
func AbsMin(x,y float32) float32 {
	if Abs(x) < Abs(y) {
		return x
	} else {
		return y
	}
}

func Mod(x,y float32) float32 {
	return float32(math.Mod(float64(x),float64(y)))
}
