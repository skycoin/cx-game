package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
)

type WindowDimensions struct {
	PhysicalWidth, PhysicalHeight float32
	VirtualWidth, VirtualHeight   float32
	Viewport                      Viewport
	Scale                         float32
}

type Viewport struct {
	X, Y, Width, Height int32
}

var currentViewport Viewport

func GetCurrentViewport() Viewport { return currentViewport }

func (v Viewport) Use() {
	//fmt.Println("Viewport: ", v)
	currentViewport = v
	gl.Viewport(v.X, v.Y, v.Width, v.Height)
}

// fits target into frame, centered with black bars
func fitCentered(virtualDims, physicalDims mgl32.Vec2) WindowDimensions {
	// "physical" dimensions describe actual window size
	// "virtual" dimensions describe scaling of both world and UI
	// physical determines resolution.
	// virtual determines how big things are.

	physicalWidth := physicalDims.X()
	physicalHeight := physicalDims.Y()
	virtualWidth := virtualDims.X()
	virtualHeight := virtualDims.Y()

	scaleToFitWidth := physicalWidth / virtualWidth
	scaleToFitHeight := physicalHeight / virtualHeight
	// scale to fit entire virtual window in physical window
	scale := cxmath.Min(scaleToFitHeight, scaleToFitWidth)

	// scale up virtual dimensions to fit in physical dimensions.
	// in case of aspect ratio mismatch, black bars will appear
	viewportWidth := int32(virtualWidth * scale)
	viewportHeight := int32(virtualHeight * scale)

	// viewport offsets
	x := (int32(physicalWidth) - viewportWidth) / 2
	y := (int32(physicalHeight) - viewportHeight) / 2

	viewport := Viewport{x, y, viewportWidth, viewportHeight}
	return WindowDimensions{
		physicalWidth, physicalHeight,
		virtualWidth, virtualHeight,
		viewport,
		scale,
	}
}

// returns a transformation matrix which maps coordinates
// on the physical window to the virtual window
func (d WindowDimensions) Transform() mgl32.Mat4 {
	// TODO
	translateToCenter := mgl32.Translate3D(
		-float32(d.PhysicalWidth)/2, -float32(d.PhysicalHeight)/2, 0,
	)
	scaleToVirtual := mgl32.Scale3D(1/d.Scale, -1/d.Scale, 1)
	// TODO cam zoom should only be applied to world coords.
	// UI coords should not be affected.
	// move this logic to the "item" package or similar.
	return mgl32.Mat4.Mul4(scaleToVirtual, translateToCenter)
}
