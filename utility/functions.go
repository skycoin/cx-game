package utility

import "math"

//clamp value
func Clamp(value, min, max float32) float32 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func boolToFloat(x bool) float32 {
	if x {
		return 1
	} else {
		return 0
	}
}

// DegToRad converts degrees to radians
func DegToRad(angle float32) float32 {
	return angle * float32(math.Pi) / 180
}

// RadToDeg converts radians to degrees
func RadToDeg(angle float32) float32 {
	return angle * 180 / float32(math.Pi)
}

func Lerp(start, finish, percent float32) float32 {
	if percent > 1 {
		return finish
	}
	return (percent * finish) + ((1 - percent) * start)
}
