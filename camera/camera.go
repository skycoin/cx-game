package camera

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/input"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
)

type zoomStatus int

var (
	//variables for smooth zooming over time
	zooming      bool    = false // flag for to know if zooming is occuring
	zoomDuration float32 = 0.6   // in seconds

	// variables for interpolation
	zoomTarget   float32 = 1 // zoom value to end on
	zoomCurrent  float32     // zoom value to start from
	zoomProgress float32     // current zoom progress (from 0 to 1)
	firstTick    bool    = true
)

type Camera struct {
	X        float32
	Y        float32
	Zoom     float32
	movSpeed float32
	window   *render.Window
	Frustum  Frustum
}

//Initiates Camera Instances given the window
func NewCamera(window *render.Window) *Camera {
	cam := Camera{
		//take X,Y pos as a center to frustrum
		X:        0,
		Y:        0,
		Zoom:     1,
		movSpeed: 5,
		window:   window,
	}

	return &cam
}

//Updates Camera Positions
func (camera *Camera) MoveCam(dTime float32) {
	camera.X += input.GetAxis(input.HORIZONTAL) * dTime * camera.movSpeed
	camera.Y += input.GetAxis(input.VERTICAL) * dTime * camera.movSpeed
	camera.UpdateFrustum()
}

func (camera *Camera) GetView() mgl32.Mat4 {
	return mgl32.Translate3D(-camera.X, -camera.Y, -camera.Zoom)
}

func (camera *Camera) GetProjectionMatrix() mgl32.Mat4 {
	left := -float32(camera.window.Width) / 2 / 32 / camera.Zoom
	right := float32(camera.window.Width) / 2 / 32 / camera.Zoom
	bottom := -float32(camera.window.Height) / 2 / 32 / camera.Zoom
	top := float32(camera.window.Height) / 2 / 32 / camera.Zoom
	projection := mgl32.Ortho(left, right, bottom, top, -1, 1000)

	return projection
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
	if !zooming {
		zooming = true
		zoomCurrent = camera.Zoom

		//find better values for better zooming
		var offset float32
		// if zoomCurrent == 1 || zoomCurrent == 2.25 {
		// 	offset = 0.75 * zoomOffset
		// } else {
		// 	offset = 1.25 * zoomOffset
		// }
		offset = 0.5 * zoomOffset

		zoomTarget = zoomCurrent + offset/2
		zoomTarget = utility.Clamp(zoomTarget, 0.5, 3)

		if zoomTarget == zoomCurrent {
			zooming = false
		}
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
	// projection := camera.GetProjectionMatrix()
	// gl.UseProgram(camera.window.Program)
	// gl.UniformMatrix4fv(gl.GetUniformLocation(camera.window.Program, gl.Str("projection\x00")), 1, false, &projection[0])
	projection := camera.GetProjectionMatrix()
	camera.window.SetProjectionMatrix(projection)
}

func (camera *Camera) Tick(dt float32) {
	// TODO optimize this later if necessary
	// always update the projection matrix in case window got resized
	camera.updateProjection()

	if firstTick {
		camera.updateProjection()
		firstTick = false
	}
	// if not zooming nothing to do here
	if !zooming {
		return
	}

	zoomProgress += dt / zoomDuration

	camera.Zoom = utility.Lerp(zoomCurrent, zoomTarget, zoomProgress)
	camera.updateProjection()

	if camera.Zoom == zoomTarget {
		zooming = false
		zoomProgress = 0
	}
}
