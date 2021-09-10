package camera

import (
	"github.com/skycoin/cx-game/cxmath"
)

var (
	//distance from center to to left/right edges
	halfWidth float32 = 24
	//distance from center to top/bottom edges
	halfHeight float32 = 16
)

func (camera *Camera) UpdateFrustum() {
	camera.Frustum.Left = int(camera.X) - int(halfWidth/camera.Zoom)
	camera.Frustum.Right = int(camera.X) + int(halfWidth/camera.Zoom)
	camera.Frustum.Top = int(camera.Y) + int(halfHeight/camera.Zoom)
	camera.Frustum.Bottom = int(camera.Y) - int(halfHeight/camera.Zoom)
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

	bottom := camera.Frustum.Bottom
	top := camera.Frustum.Top
	left := camera.Frustum.Left
	right := camera.Frustum.Right
	posX := int(position.X)
	posY := int(position.Y)

	if posY+radius >= bottom && posY-radius <= top {
		if posX+radius >= left && posX+radius <= right {
			return true
		}
		if left < 0 {
			left += int(camera.PlanetWidth)
			if posX > left {
				return true
			}
		}
		if right > int(camera.PlanetWidth) {
			right -= int(camera.PlanetWidth)
			if posX < right {
				return true
			}
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
