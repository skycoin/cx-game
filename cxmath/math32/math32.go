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

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Clamp(x, min, max float32) float32 {
	if x < min { return min }
	if x > max { return max }
	return x
}

// non negative solution to x % d
func PositiveModulo(x, b float32) float32 {
	x = Mod(x,b)
	if x >= 0 {
		return x
	}
	return x + b
}
