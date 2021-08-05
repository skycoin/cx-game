package physics

//https://www.youtube.com/watch?v=PlT44xr0iW0
type CollisionInfo struct {
	Above bool
	Below bool
	Left  bool
	Right bool
}

func (c *CollisionInfo) Reset() {
	c.Above, c.Below, c.Left, c.Right = false, false, false, false
}

func (c CollisionInfo) Horizontal() bool {
	return c.Left || c.Right
}
