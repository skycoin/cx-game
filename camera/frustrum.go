package camera

type Frustrum struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

var (
	baseWidth  float32 = 32
	baseHeight float32 = 32
	debugUseSmallerFrustrum bool = false
)

type FrustumSize struct {
	width float32
	height float32
}

func (cam *Camera) ToggleDebugUseSmallerFrustrum() {
	debugUseSmallerFrustrum = !debugUseSmallerFrustrum
	cam.UpdateFrustrum()
}

func (cam *Camera) getFrustrumSize() FrustumSize {
	if debugUseSmallerFrustrum {
		return FrustumSize { 10, 10 }
	} else {
		return FrustumSize { baseWidth, baseHeight }
	}
}

func (cam *Camera) UpdateFrustrum() {
	baseSize := cam.getFrustrumSize()

	cam.Frustrum.Left = int(cam.X - baseSize.width/cam.Zoom)
	cam.Frustrum.Right = int(cam.X + baseSize.width/cam.Zoom)
	cam.Frustrum.Top = int(cam.Y + baseSize.height/cam.Zoom)
	cam.Frustrum.Bottom = int(cam.Y - baseSize.height/cam.Zoom)
}
