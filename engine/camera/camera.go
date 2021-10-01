package camera

import (
	"fmt"
	"log"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/effects"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/render"
)

var (
	//variables for smooth zooming over time
	// zooming      bool    = false // flag for to know if zooming is occuring
	// zoomDuration float32 = 0.4   // in seconds
	// 	zoomProgress float32

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
	//main camera
	X          float32
	Y          float32
	Zoom       Zoom
	Frustum    cxmath.Frustum
	isOnTarget bool
	//player position
	PlayerX float32
	PlayerY float32
	// PlayerZoom float32
	//target position
	TargetX float32
	TargetY float32
	// TargetZoom float32
	//rest of the code
	Vel      mgl32.Vec2
	movSpeed float32
	window   *render.Window
	//focus_area  focusArea
	freeCam     bool
	PlanetWidth float32
	shake       *effects.ShakeStruct
	shockwave   *effects.Shockwave
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
	// size := mgl32.Vec2{3, 5}
	// xPos := float32(0)
	// yPos := float32(0)
	cam := Camera{
		//take X,Y pos as a center to frustrum
		// X:        xPos,
		// Y:        yPos,
		Vel:      mgl32.Vec2{0, 0},
		Zoom:     NewZoom(),
		movSpeed: 5,
		window:   window,
		// focus_area: focusArea{
		// 	center: mgl32.Vec2{xPos, yPos},
		// 	left:   xPos - size.X()/2,
		// 	right:  xPos + size.X()/2,
		// 	top:    yPos + size.Y()/2,
		// 	bottom: yPos - size.Y()/2,
		// },
		freeCam:   false,
		shake:     effects.NewShakeStruct(60, 1),
		shockwave: effects.NewShockwave(),
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
	return mgl32.Translate3D(-camera.X, -camera.Y, -camera.Zoom.Get())
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
	camera.PlayerX = x
	camera.PlayerY = y
	// camera.updateFocusArea(x, y)
}

var CameraSnapped bool

// update focus area to include (x,y)
// func (camera *Camera) updateFocusArea(xPos, yPos float32) {
// 	x, y := ApplySnapping(xPos, yPos)
// 	/*
// 		no camera snap
// 		snap camera center to nearest 1.0f / 32 pixel
// 		snap camera center to nearest 1.0f / 32 + 1.0f / 64 pixel
// 		snap camera center to nearest 1.0f / 16 pixel
// 		snap camera center to nearest 1.0f / 16 + 1.0f / 32 pixel
// 	*/
// 	modular := cxmath.NewModular(camera.PlanetWidth)
// 	var shiftX, shiftY float32
// 	if modular.IsLeft(x, camera.focus_area.left) {
// 		shiftX = modular.Disp(camera.focus_area.left, x)
// 	} else if modular.IsRight(x, camera.focus_area.right) {
// 		shiftX = modular.Disp(camera.focus_area.right, x)
// 	}
// 	if y < camera.focus_area.bottom {
// 		shiftY = y - camera.focus_area.bottom
// 	} else if y > camera.focus_area.top {
// 		shiftY = y - camera.focus_area.top
// 	}
// 	camera.focus_area.left += shiftX
// 	camera.focus_area.right += shiftX
// 	camera.focus_area.bottom += shiftY
// 	camera.focus_area.top += shiftY
// 	camera.focus_area.center = mgl32.Vec2{
// 		(camera.focus_area.left + camera.focus_area.right) / 2,
// 		(camera.focus_area.top + camera.focus_area.bottom) / 2,
// 	}

// 	camera.Vel[0] = math32.Mod((camera.focus_area.center.X() - camera.X), 100)
// 	camera.Vel[1] = camera.focus_area.center.Y() - camera.Y

// 	camera.X = math32.
// 		PositiveModulo(camera.focus_area.center.X(), camera.PlanetWidth)
// 	camera.Y = camera.focus_area.center.Y()

// }

//sets camera for target position
func (camera *Camera) SetCameraPositionTarget(x, y float32) {
	camera.TargetX = x
	camera.TargetY = y
}

//zooms on target
// func (camera *Camera) SetCameraZoomTarget(zoomOffset float32) {
// 	camera.TargetZoom =
// }

//zooms on current position
func (camera *Camera) SetCameraZoomPosition(zoomOffset float32) {
	if zoomOffset > 0 {
		camera.Zoom.Up()
		return
	}
	camera.Zoom.Down()
}

func (camera Camera) GetTransform() mgl32.Mat4 {
	xoff := camera.X
	yoff := camera.Y

	if camera.shake.IsShaking {
		xoff += camera.shake.Amplitude()
		yoff += camera.shake.Amplitude()
	}
	translate := mgl32.Translate3D(xoff, yoff, 0)

	// fmt.Println(camera.Zoom)
	scale := mgl32.Scale3D(1/camera.Zoom.Get(), 1/camera.Zoom.Get(), 1)
	return translate.Mul4(scale)
}

func (camera *Camera) Tick(dt float32) {
	if !camera.isOnTarget {
		camera.X = camera.PlayerX
		camera.Y = camera.PlayerY
	} else {
		camera.X = camera.TargetX
		camera.Y = camera.TargetY
	}
	camera.UpdateFrustum()
	camera.Zoom.Tick(dt)
	camera.shake.Update(dt)

	camera.shockwave.SetZoom(camera.Zoom.Get())
	camera.shockwave.Update(dt)
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
	switch camera.Zoom.value {
	case 1:
		camera.Zoom.Set(0.5)
		message = "2 pixels per pixel"
	case 0.5:
		camera.Zoom.Set(1.0 / 3.0)
		message = "3 pixels per pixel"
	case 1.0 / 3.0:
		camera.Zoom.Set(1.0 / 4.0)
		message = "4 pixels per pixel"
	default:
		camera.Zoom.Set(1)
		message = "1 pixel per pixels"
	}
	fmt.Printf("Current zoom: %v\n", message)
}

func (camera *Camera) Pos() mgl32.Vec2 {
	return mgl32.Vec2{camera.X, camera.Y}
}

func (camera *Camera) SwitchToTarget() {
	camera.isOnTarget = true
}
func (camera *Camera) SwitchToPlayer() {
	camera.isOnTarget = false
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

func (cam *Camera) Shake() {
	cam.shake.Start(60, 0.8)
}

//start shockwave at
func (cam *Camera) Shockwave(x, y float32) {

	screenPos := cam.ScreenPos(mgl32.Vec2{x, y}).Add(mgl32.Vec2{0.5, 0.5})

	cam.shockwave.SetCenter(screenPos.X(), screenPos.Y())
	cam.shockwave.Start()
}

func (cam *Camera) ShockwaveVec(pos cxmath.Vec2) {
	cam.Shockwave(pos.X, pos.Y)
}

//convert cam position into NDC (-1,1) coords
func (cam *Camera) ScreenPos(pos mgl32.Vec2) mgl32.Vec2 {

	converted := render.Projection.Mul4(
		cam.GetViewMatrix()).Mul4x1(
		pos.Vec4(0, 1),
	)
	return converted.Vec2()
}
