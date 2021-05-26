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
