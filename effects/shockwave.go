package effects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
)

//cant const mgl32.Vec2
var (
	center_default            = mgl32.Vec2{0.5, 0.5}
	force_default     float32 = 0.2
	size_default      float32 = 0.05
	size_max_default  float32 = 0.23
	thickness_default float32 = 0.0
	duration_default  float32 = 0.5
)

type Shockwave struct {
	center    mgl32.Vec2
	force     float32
	size      float32
	thickness float32
	duration  float32
	timeLeft  float32
	IsActive  bool
	zoom      float32
}

func NewShockwave() *Shockwave {
	return &Shockwave{
		center:    center_default,
		force:     force_default,
		size:      size_default,
		thickness: thickness_default,
		duration:  duration_default,
	}
}

func (s *Shockwave) Update(dt float32) {
	// fmt.Printf("IS Active: %v, duration: %v,  timeLeft: %v, size: %v, center: %v\n",
	// 	s.IsActive,
	// 	s.duration,
	// 	s.timeLeft,
	// 	s.size,
	// 	s.center,
	// )
	if !s.IsActive {
		return
	}

	s.timeLeft -= dt

	if s.timeLeft < 0 {
		s.IsActive = false
		render.DisableShockwave()
	}

	s.CalculateSize()

	render.SetShockwaveCenter(s.center.X(), s.center.Y())
	render.SetShockwaveSize(s.size * s.zoom)
	render.SetShockwaveForce(s.force * s.zoom)
	render.SetShockwaveThickness(s.thickness * s.zoom)
}

//calculate size from duration, size dissappears at 0.8
func (s *Shockwave) CalculateSize() {
	// fmt.Println(s.timeLeft)
	s.size = size_default + size_max_default*cxmath.Interpolate(0, 1, (s.duration-s.timeLeft)/s.duration, cxmath.EASEOUTSINE)

}

func (s *Shockwave) SetDuration(sec float32) {
	s.duration = sec
	if !s.IsActive {
		s.timeLeft = s.duration
	}
}

func (s *Shockwave) SetCenter(x, y float32) {
	s.center = mgl32.Vec2{x, y}
}

func (s *Shockwave) Start() {
	if s.IsActive {
		return
	}
	s.timeLeft = s.duration
	s.IsActive = true

	render.EnableShockwave()
}

func (s *Shockwave) SetZoom(zoom float32) {
	s.zoom = zoom
}
