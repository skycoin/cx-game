package spine

type Animation struct {
	Name      string
	Duration  float32
	Timelines []Timeline
}

type Timeline interface {
	Duration() float32
	Apply(skeleton *Skeleton, time float32, alpha float32)
}

func (anim *Animation) Apply(skeleton *Skeleton, time float32, loop bool) {
	anim.Mix(skeleton, time, loop, 1)
}

func (anim *Animation) Mix(skeleton *Skeleton, time float32, loop bool, alpha float32) {
	if loop && anim.Duration > 0 {
		time = mod(time, anim.Duration)
	}
	for _, timeline := range anim.Timelines {
		timeline.Apply(skeleton, time, alpha)
	}
}
