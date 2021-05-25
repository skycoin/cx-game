package camera

type Frustum struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

var (
	cameraCurrent Frustum
	// cameraTarget Frustrum

	//distance from center to to left/right edges
	halfWidth float32 = 16
	//distance from center to top/bottom edges
	halfHeight float32 = 16
)

func (camera *Camera) UpdateFrustrum() {
	camera.Frustum.Left = int(camera.X - halfWidth/camera.Zoom)
	camera.Frustum.Right = int(camera.X + halfWidth/camera.Zoom)
	camera.Frustum.Top = int(camera.Y + halfHeight/camera.Zoom)
	camera.Frustum.Bottom = int(camera.Y - halfHeight/camera.Zoom)
}
