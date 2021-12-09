package mathi

func Min(x,y int) int {
	if x<y { return x } else { return y }
}

func Max(x,y int) int {
	if x>y { return x } else { return y }
}

func Clamp(x, min, max int) int {
	if x < min { return min }
	if x > max { return max }
	return x
}

func Abs(x int) int {
	if x<0 { return -x } else { return x }
}
