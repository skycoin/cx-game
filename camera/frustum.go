package camera

type Frustum struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

var (
	//distance from center to to left/right edges
	halfWidth float32 = 16
	//distance from center to top/bottom edges
	halfHeight float32 = 16
	//margin
	margin = 3
)

func (camera *Camera) UpdateFrustum() {
	camera.Frustum.Left = int(camera.X) - margin - int(halfWidth)
	camera.Frustum.Right = int(camera.X) + margin + int(halfWidth)
	camera.Frustum.Top = int(camera.Y) + margin + int(halfHeight)
	camera.Frustum.Bottom = int(camera.Y) - margin - int(halfHeight)
}

func (camera *Camera) IsInBounds(x, y int) bool {
	if x >= camera.Frustum.Left &&
		x <= camera.Frustum.Right &&
		y >= camera.Frustum.Bottom &&
		y <= camera.Frustum.Top {
		return true
	}
	return false
}

func (camera *Camera) IsInBoundsF(x, y float32) bool {
	return camera.IsInBounds(int(x), int(y))
}
