package camera

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	X        float32
	Y        float32
	Zoom     float32
	movSpeed float32
}

func NewCamera() *Camera {
	cam := Camera{
		X:        0,
		Y:        0,
		Zoom:     0,
		movSpeed: 0,
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

func (camera *Camera) GetTransform() mgl32.Mat4 {
	return mgl32.Translate3D(camera.X, camera.Y, camera.Zoom)
}
