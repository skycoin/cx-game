package camera

import (
	"github.com/skycoin/cx-game/cxmath"
)

var (
	// variables for interpolation

	zoomLevels = []float32{
		0.5, 1, 2, 4, 8, 16,
	}
)

type Zoom struct {
	value        float32
	isZooming    bool
	zoomProgress float32
	duration     float32
	// zoomLevels       []float32
	zoomCurrent      float32
	zoomNext         float32
	currentZoomIndex int
}

func NewZoom() Zoom {
	zoom := Zoom{
		duration:         0.5,
		value:            1,
		zoomCurrent:      1,
		zoomNext:         1,
		currentZoomIndex: 1,
	}
	return zoom
}

func (z *Zoom) Tick(dt float32) {
	if !z.isZooming {
		return
	}

	z.zoomProgress += dt / z.duration
	if z.zoomProgress > 1 {
		z.isZooming = false
		z.value = z.zoomNext
		z.zoomCurrent = z.zoomNext
		z.zoomProgress = 0
	} else {
		z.value = cxmath.Lerp(z.zoomCurrent, z.zoomNext, z.zoomProgress)
	}
}

//instant
func (z *Zoom) Set(value float32) {
	z.zoomNext = value
	z.zoomCurrent = value

}

func (z Zoom) Get() float32 {
	return z.value
}

func (z *Zoom) SetDuration(duration float32) {
	z.duration = duration
}

func (z Zoom) GetFrustum() float32 {
	if z.isZooming {
		if z.zoomCurrent > z.zoomNext {
			return z.zoomNext
		}
		return z.zoomCurrent
	}

	return z.value

}

func (z *Zoom) Up() {
	if z.currentZoomIndex == len(zoomLevels)-1 || z.isZooming {
		return
	}
	z.currentZoomIndex += 1
	z.zoomNext = zoomLevels[z.currentZoomIndex]
	z.isZooming = true
}

func (z *Zoom) Down() {
	if z.currentZoomIndex == 0 || z.isZooming {
		return
	}
	z.currentZoomIndex -= 1
	z.zoomNext = zoomLevels[z.currentZoomIndex]
	z.isZooming = true
}
