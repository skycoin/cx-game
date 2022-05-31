package spine

import (
	"math"
)

func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}
func atan2(y, x float32) float32 { return float32(math.Atan2(float64(y), float64(x))) }
func mod(x, y float32) float32   { return float32(math.Mod(float64(x), float64(y))) }
func sqrt(v float32) float32     { return float32(math.Sqrt(float64(v))) }
func sin(v float32) float32      { return float32(math.Sin(float64(v))) }
func cos(v float32) float32      { return float32(math.Cos(float64(v))) }
func sincos(v float32) (float32, float32) {
	sn, cs := math.Sincos(float64(v))
	return float32(sn), float32(cs)
}

func clamp01(v float32) float32 {
	if v <= 0 {
		return 0
	} else if v >= 1 {
		return 1
	}
	return v
}

func lerp(a, b float32, p float32) float32 { return a*(1-p) + b*p }

func lerpAngle(a, b, p float32) float32 {
	delta := b - a
	for delta > math.Pi {
		delta -= 2 * math.Pi
	}
	for delta < -math.Pi {
		delta += 2 * math.Pi
	}
	return a + delta*p
}

func lerpVector(a, b Vector, p float32) Vector {
	return Vector{
		lerp(a.X, b.X, p),
		lerp(a.Y, b.Y, p),
	}
}

func lerpAngleVector(a, b Vector, p float32) Vector {
	return Vector{
		lerpAngle(a.X, b.X, p),
		lerpAngle(a.Y, b.Y, p),
	}
}
