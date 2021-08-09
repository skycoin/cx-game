package camera

import (
	"github.com/skycoin/cx-game/cxmath"
)

var (
	//distance from center to to left/right edges
	halfWidth float32 = 16
	//distance from center to top/bottom edges
	halfHeight float32 = 16
	//margin
	margin = 3
)

func (camera *Camera) UpdateFrustum() {
	camera.Frustum.Left = int(camera.X) - margin - int(halfWidth)
	camera.Frustum.Right = int(camera.X) + margin + int(halfWidth)
	camera.Frustum.Top = int(camera.Y) + margin + int(halfHeight)
	camera.Frustum.Bottom = int(camera.Y) - margin - int(halfHeight)
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

func (camera *Camera) IsInBoundsRadius(position cxmath.Vec2, radius int) bool {
	//https://imgur.com/a/HsuerD3
	// (x - X )^2 + (y - Y) ^2 and if its more than R^2

	//check left and right
	if position.X >= float32(camera.Frustum.Bottom) && position.Y <= float32(camera.Frustum.Top) {
		if position.X+float32(radius) >= float32(camera.Frustum.Left) || position.X-float32(radius) <= float32(camera.Frustum.Right) {
			return true
		}
		//up and down
	} else if position.X >= float32(camera.Frustum.Left) && position.X <= float32(camera.Frustum.Right) {
		if position.Y+float32(radius) >= float32(camera.Frustum.Bottom) || position.Y-float32(radius) <= float32(camera.Frustum.Top) {
			return true
		}
	}
	return false
}

func (camera *Camera) TilesInView() []cxmath.Vec2i {
	positions := []cxmath.Vec2i{}
	for y := camera.Frustum.Bottom; y <= camera.Frustum.Top; y++ {
		for x := camera.Frustum.Left; x <= camera.Frustum.Right; x++ {
			positions = append(positions, cxmath.Vec2i{int32(x), int32(y)})
		}
	}
	return positions
}
