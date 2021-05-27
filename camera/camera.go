package camera

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
)

type zoomStatus int

var (
	zooming     bool    = false
	zoomTarget  float32 = 1
	zoomCurrent float32
	zoomPercent float32
	zoomSpeed   float32 = 1.3
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
	camera.UpdateFrustum()
}

//moves and/or zooms  camera
func (camera *Camera) GetView() mgl32.Mat4 {
	return mgl32.Translate3D(-camera.X, -camera.Y, camera.Zoom)
}

func (camera *Camera) SetCameraCenter() {
	camera.X = float32(camera.window.Width) / 2
	camera.Y = float32(camera.window.Height) / 2
	camera.UpdateFrustum()
}

//sets camera for current position
func (camera *Camera) SetCameraPosition(x, y float32) {
	camera.X = x
	camera.Y = y
	camera.UpdateFrustum()
}

//sets camera for target position
func (camera *Camera) SetCameraPositionTarget(x, y float32) {
	camera.SetCameraPosition(x, y)
	camera.UpdateFrustum()
}

//zooms on target
func (camera *Camera) SetCameraZoomTarget(zoomOffset float32) {
	camera.SetCameraZoomPosition(zoomOffset)
}

//zooms on current position

func (camera *Camera) SetCameraZoomPosition(zoomOffset float32) {
	// camera.Zoom += zoomOffset
	// camera.Zoom = utility.Clamp(camera.Zoom, 1, 3)
	// camera.updateProjection()
	if !zooming {
		zooming = true
		zoomCurrent = camera.Zoom
		zoomTarget = zoomCurrent + zoomOffset
		zoomTarget = utility.Clamp(zoomTarget, 1, 3)
	}
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

func (camera *Camera) updateProjection() {
	left := -float32(camera.window.Width) / 2 / 32 / camera.Zoom
	right := float32(camera.window.Width) / 2 / 32 / camera.Zoom
	bottom := -float32(camera.window.Height) / 2 / 32 / camera.Zoom
	top := float32(camera.window.Height) / 2 / 32 / camera.Zoom
	projection := mgl32.Ortho(left, right, bottom, top, -1, 1000)

	gl.UseProgram(camera.window.Program)
	gl.UniformMatrix4fv(gl.GetUniformLocation(camera.window.Program, gl.Str("projection\x00")), 1, false, &projection[0])
}

func (camera *Camera) Tick(dt float32) {
	// if not zooming nothing to do here
	if !zooming {
		return
	}
	zoomPercent += dt * zoomSpeed

	camera.Zoom = utility.Lerp(zoomCurrent, zoomTarget, zoomPercent)
	camera.updateProjection()

	// fmt.Println(camera.Zoom, "    ", zoomTarget)
	if camera.Zoom == zoomTarget {
		zooming = false
		zoomPercent = 0
	}
}
