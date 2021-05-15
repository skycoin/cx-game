package camera

type Frustrum struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

var (
	cameraCurrent Frustrum
	// cameraTarget Frustrum
	baseWidth  float32 = 32
	baseHeight float32 = 32
)

func UpdateFrustrum(x, y, zoom float32) {
	cameraCurrent.Left = int(x - baseWidth/zoom)
	cameraCurrent.Right = int(x + baseWidth/zoom)
	cameraCurrent.Top = int(y + baseHeight/zoom)
	cameraCurrent.Bottom = int(y - baseHeight/zoom)
}
