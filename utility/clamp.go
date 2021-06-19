package utility

func ClampF(value, min, max float32) float32 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
func ClampI(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
