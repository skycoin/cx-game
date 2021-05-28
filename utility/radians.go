package utility

import "math"

// DegToRad converts degrees to radians
func DegToRad(angle float32) float32 {
	return angle * float32(math.Pi) / 180
}

// RadToDeg converts radians to degrees
func RadToDeg(angle float32) float32 {
	return angle * 180 / float32(math.Pi)
}
