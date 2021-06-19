package utility

//interpolate from start to finish in duration time
//		alpha - 123
func Lerp(start, target, alpha float32) float32 {
	alpha = ClampF(alpha, 0, 1)
	return (alpha * target) + ((1 - alpha) * start)
}

func SmoothStep(start, target float32, x float32) float32 {
	x = ClampF((x-start)/(target-start), 0, 1)

	return x * x * (3 - 2*x)
}
