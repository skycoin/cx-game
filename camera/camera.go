package camera

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
)

type Camera struct {
	X        float32
	Y        float32
	Zoom     float32
	movSpeed float32
	window   *render.Window
	Frustum  Frustum
}

func NewCamera(window *render.Window) *Camera {
	cam := Camera{
		//take X,Y pos as a center to frustrum
		X:        0,
		Y:        0,
		Zoom:     1.0,
		movSpeed: 5,
		window:   window,
	}

	return &cam
}

func (camera *Camera) MoveCam(x, y float32, dTime float32) {
	camera.X += x * dTime * camera.movSpeed
	camera.Y += y * dTime * camera.movSpeed
	// camera.Zoom += z * dTime * camera.movSpeed
	camera.UpdateFrustrum()
}

//moves and/or zooms  camera
func (camera *Camera) GetView() mgl32.Mat4 {
	return mgl32.Translate3D(-camera.X, -camera.Y, camera.Zoom)
}

func (camera *Camera) SetCameraCenter() {
	camera.X = float32(camera.window.Width) / 2
	camera.Y = float32(camera.window.Height) / 2
	camera.UpdateFrustrum()
}

//sets camera for current position
func (camera *Camera) SetCameraPosition(x, y float32) {
	camera.X = x
	camera.Y = y
	camera.UpdateFrustrum()
}

//sets camera for target position
func (camera *Camera) SetCameraPositionTarget(x, y float32) {
	camera.SetCameraPosition(x, y)
	camera.UpdateFrustrum()
}

//zooms on target
func (camera *Camera) SetCameraZoomTarget(zoom float32) {
	camera.SetCameraZoomPosition(zoom)
}

//zooms on current position
func (camera *Camera) SetCameraZoomPosition(zoom float32) {
	camera.Zoom = utility.Clamp(camera.Zoom+zoom, 1, 3)
	camera.UpdateFrustrum()
}

func (camera *Camera) DrawLines(
	lines []float32, color []float32, ctx render.Context,
) {
	camCtx := ctx.PushView(camera.GetView())
	camera.window.DrawLines(lines, color, camCtx)
}

func (camera Camera) GetTransform() mgl32.Mat4 {
	return mgl32.Translate3D(camera.X, camera.Y, 0)
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
