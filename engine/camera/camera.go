package camera

import (
	"fmt"
	"log"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/cxmath/mathi"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/render"
)

var (
	//variables for smooth zooming over time
	zooming      bool    = false // flag for to know if zooming is occuring
	zoomDuration float32 = 0.4   // in seconds
	zoomProgress float32

	// variables for interpolation
	zoomTarget  float32 = 1 // zoom value to end on
	zoomCurrent float32     // zoom value to start from
	// current zoom progress (from 0 to 1)

	currentZoomIndex int = 1

	zoomLevels = []float32{
		0.5, 1, 2, 4, 8, 16,
	}

	currentSnapping = NO_SNAPPING
)

const (
	//camera snapping constants
	NO_SNAPPING = iota
	// 1.0/32
	NEAREST_1
	// 1.0/32 + 1.0/64
	NEAREST_2
	// 1.0/16
	NEAREST_3
	// 1.0/16 + 1.0/32
	NEAREST_4
)

type Camera struct {
	X           float32
	Y           float32
	Vel         mgl32.Vec2
	Zoom        float32
	movSpeed    float32
	window      *render.Window
	Frustum     cxmath.Frustum
	focus_area  focusArea
	freeCam     bool
	PlanetWidth float32
}

type focusArea struct {
	center mgl32.Vec2
	left   float32
	right  float32
	top    float32
	bottom float32
}

//Initiates Camera Instances given the window
func NewCamera(window *render.Window) *Camera {
	size := mgl32.Vec2{3, 5}
	xPos := float32(0)
	yPos := float32(0)
	cam := Camera{
		//take X,Y pos as a center to frustrum
		X:        xPos,
		Y:        yPos,
		Vel:      mgl32.Vec2{0, 0},
		Zoom:     1,
		movSpeed: 5,
		window:   window,
		focus_area: focusArea{
			center: mgl32.Vec2{xPos, yPos},
			left:   xPos - size.X()/2,
			right:  xPos + size.X()/2,
			top:    yPos + size.Y()/2,
			bottom: yPos - size.Y()/2,
		},
		freeCam: false,
	}
	return &cam
}

//Updates Camera Positions
func (camera *Camera) MoveCam(dTime float32) {
	if input.GetInputContext() != input.FREECAM {
		return
	}
	camera.X += input.GetAxis(input.HORIZONTAL) * dTime * camera.movSpeed
	camera.Y += input.GetAxis(input.VERTICAL) * dTime * camera.movSpeed
	camera.UpdateFrustum()
}

func (camera *Camera) GetView() mgl32.Mat4 {
	return mgl32.Translate3D(-camera.X, -camera.Y, -camera.Zoom)
}

func (camera *Camera) GetProjectionMatrix() mgl32.Mat4 {
	ppt := float32(constants.PIXELS_PER_TILE)
	left := -float32(camera.window.Width) / 2 / ppt
	right := float32(camera.window.Width) / 2 / ppt
	bottom := -float32(camera.window.Height) / 2 / ppt
	top := float32(camera.window.Height) / 2 / ppt
	projection := mgl32.Ortho(left, right, bottom, top, -1000, 1000)

	return projection
}

func (camera *Camera) GetViewMatrix() mgl32.Mat4 {
	return camera.GetTransform().Inv()
	// return camera.GetTransform()
}

func (camera *Camera) SetCameraCenter() {
	camera.X = float32(camera.window.Width) / 2
	camera.Y = float32(camera.window.Height) / 2
	camera.UpdateFrustum()
}

//sets camera for current position
func (camera *Camera) SetCameraPosition(x, y float32) {
	camera.updateFocusArea(x, y)
	camera.UpdateFrustum()
}

var CameraSnapped bool

// update focus area to include (x,y)
func (camera *Camera) updateFocusArea(xPos, yPos float32) {
	x, y := ApplySnapping(xPos, yPos)
	/*
		no camera snap
		snap camera center to nearest 1.0f / 32 pixel
		snap camera center to nearest 1.0f / 32 + 1.0f / 64 pixel
		snap camera center to nearest 1.0f / 16 pixel
		snap camera center to nearest 1.0f / 16 + 1.0f / 32 pixel
	*/
	modular := cxmath.NewModular(camera.PlanetWidth)
	var shiftX, shiftY float32
	if modular.IsLeft(x, camera.focus_area.left) {
		shiftX = modular.Disp(camera.focus_area.left, x)
	} else if modular.IsRight(x, camera.focus_area.right) {
		shiftX = modular.Disp(camera.focus_area.right, x)
	}
	if y < camera.focus_area.bottom {
		shiftY = y - camera.focus_area.bottom
	} else if y > camera.focus_area.top {
		shiftY = y - camera.focus_area.top
	}
	camera.focus_area.left += shiftX
	camera.focus_area.right += shiftX
	camera.focus_area.bottom += shiftY
	camera.focus_area.top += shiftY
	camera.focus_area.center = mgl32.Vec2{
		(camera.focus_area.left + camera.focus_area.right) / 2,
		(camera.focus_area.top + camera.focus_area.bottom) / 2,
	}

	camera.Vel[0] = math32.Mod((camera.focus_area.center.X() - camera.X), 100)
	camera.Vel[1] = camera.focus_area.center.Y() - camera.Y

	camera.X = math32.
		PositiveModulo(camera.focus_area.center.X(), camera.PlanetWidth)
	camera.Y = camera.focus_area.center.Y()

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
		zoomCurrent = zoomLevels[currentZoomIndex]

		currentZoomIndex = mathi.Clamp(
			currentZoomIndex+int(zoomOffset),
			0, len(zoomLevels)-1,
		)
		nextZoomIndex := currentZoomIndex

		zoomTarget = zoomLevels[nextZoomIndex]
		if zoomTarget == zoomCurrent {
			zooming = false
		}
	}
}

func (camera Camera) GetTransform() mgl32.Mat4 {
	translate := mgl32.Translate3D(camera.X, camera.Y, 0)
	// fmt.Println(camera.Zoom)
	scale := mgl32.Scale3D(1/camera.Zoom, 1/camera.Zoom, 1)
	return translate.Mul4(scale)
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

	// if firstTick {
	// 	camera.updateProjection()
	// 	firstTick = false
	// }
	// if not zooming nothing to do here
	if !zooming {
		return
	}

	zoomProgress += dt / zoomDuration

	camera.Zoom = cxmath.Lerp(zoomCurrent, zoomTarget, zoomProgress)
	camera.updateProjection()

	if camera.Zoom == zoomTarget {
		zooming = false
		zoomProgress = 0
	}
}

func (camera *Camera) IsFreeCam() bool {
	return camera.freeCam
}

func (camera *Camera) TurnOnFreeCam() {
	camera.freeCam = true
	input.SetInputContext(input.FREECAM)
}
func (camera *Camera) TurnOffFreeCam() {
	camera.freeCam = false
	input.SetInputContext(input.GAME)
}
func (camera *Camera) CycleZoom() {
	var message string
	switch camera.Zoom {
	case 1:
		camera.Zoom = 0.5
		message = "2 pixels per pixel"
	case 0.5:
		camera.Zoom = 1.0 / 3.0
		message = "3 pixels per pixel"
	case 1.0 / 3.0:
		camera.Zoom = 1.0 / 4.0
		message = "4 pixels per pixel"
	default:
		camera.Zoom = 1
		message = "1 pixel per pixels"
	}
	fmt.Printf("Current zoom: %v\n", message)
}

func (c *Camera) Pos() mgl32.Vec2 {
	return mgl32.Vec2{c.X, c.Y}
}

func CycleSnap() {
	currentSnapping = (currentSnapping + 1) % 5
	DebugSnapping(currentSnapping)
}

func ApplySnapping(x, y float32) (float32, float32) {
	var newX, newY float32
	switch currentSnapping {
	case NO_SNAPPING:
		return x, y
	case NEAREST_1:
		newX, newY = math32.Round(x*32)/32, math32.Round(y*32)/32
	case NEAREST_2:
		newX, newY = math32.Round(x*32)/32+1.0/64, math32.Round(y*32)/32
	case NEAREST_3:
		newX, newY = math32.Round(x*16)/16, math32.Round(y*16)/16
	case NEAREST_4:
		newX, newY = math32.Round(x*16)/16+1.0/32, math32.Round(y*16)/16
	default:
		log.Fatalf("Camera snapping error")
	}
	return newX, newY

}

func DebugSnapping(option int) {
	if option == NO_SNAPPING {
		fmt.Println("[CAMERA SNAP] No snapping")
	} else if option == NEAREST_1 {
		fmt.Println("[CAMERA SNAP] Nearest 1.0/32")
	} else if option == NEAREST_2 {
		fmt.Println("[CAMERA SNAP] Nearest 1.0/32 + 1.0/64 offset")
	} else if option == NEAREST_3 {
		fmt.Println("[CAMERA SNAP] Nearest 1.0/16")
	} else if option == NEAREST_4 {
		fmt.Println("[CAMERA SNAP] Nearest 1.0/16 + 1.0/32 offset")
	} else {
		log.Fatal("No such camera snapping option, crashing")
	}
}
