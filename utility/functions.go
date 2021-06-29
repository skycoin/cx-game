package utility

import "math"

//clamp value

func BoolToFloat(x bool) float32 {
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

func RemoveElementF(slice []float32, index int) []float32 {
	if index < 0 || len(slice) >= index {
		return slice
	}
	newSlice := append(slice[:index], slice[index+1:]...)
	return newSlice
}

func RemoveElementI(slice []int, index int) []int {
	if index < 0 || len(slice) >= index {
		return slice
	}
	newSlice := append(slice[:index], slice[index+1:]...)
	return newSlice
}
