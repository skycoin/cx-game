package utility

func Lerp(start, finish, percent float32) float32 {
	if percent > 1 {
		return finish
	}
	return (percent * finish) + ((1 - percent) * start)
}
