package camera

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	X        float32
	Y        float32
	Zoom     float32
	movSpeed float32
	window   *Window
}

func NewCamera(window *Window) *Camera {
	cam := Camera{
		X:        0,
		Y:        0,
		Zoom:     1.0,
		movSpeed: 0,
		window:   window,
	}

	return &cam
}

func (camera *Camera) MoveCam(x, y, z float32, dTime float32) {
	camera.X += x * dTime * camera.movSpeed
	camera.Y += y * dTime * camera.movSpeed
	camera.Zoom += z * dTime * camera.movSpeed
	if camera.Zoom > 0 {
		camera.Zoom = 0
	}
}

//moves and/or zooms  camera
func (camera *Camera) GetTransform() mgl32.Mat4 {
	return mgl32.Translate3D(camera.X, camera.Y, camera.Zoom)
}

func (camera *Camera) SetCameraCenter() {
	camera.X = float32(camera.window.X) / 2
	camera.Y = float32(camera.window.Y) / 2
	UpdateFrustrum(camera.X, camera.Y, camera.Zoom)
}

//sets camera for current position
func (camera *Camera) SetCameraPosition(x, y float32) {
	camera.X = x
	camera.Y = y
	UpdateFrustrum(x, y, camera.Zoom)
}

//sets camera for target position
func (camera *Camera) SetCameraPositionTarget(x, y float32) {
	camera.SetCameraPosition(x, y)
}

//zooms on target
func (camera *Camera) SetCameraZoomTarget(zoom float32) {
	camera.SetCameraZoomPosition(zoom)
}

//zooms on current position
func (camera *Camera) SetCameraZoomPosition(zoom float32) {
	camera.Zoom = zoom
	UpdateFrustrum(camera.X, camera.Y, zoom)
}
