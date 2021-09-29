package camera

import (
	"math"
	"math/rand"
)

const (
	AMPLITUDE = 0.3
)

type ShakeStruct struct {
	frequency float32
	duration  float32
	samples   []float32
	t         float32
	IsShaking bool
}

func NewShakeStruct(freq, dur float32) *ShakeStruct {
	s := ShakeStruct{}
	s.frequency = freq
	s.duration = dur
	return &s
}

func (s *ShakeStruct) Start(freq, dur float32) {
	if freq != 0 {
		s.frequency = freq
	}
	if dur != 0 {
		s.duration = dur
	}

	sampleCount := int((s.duration) * s.frequency)

	// fmt.Println("sampleCount: ", sampleCount)
	// Populate the samples array with randomized values between -1.0 and 1.0
	s.samples = make([]float32, sampleCount)
	for i := range s.samples {
		s.samples[i] = rand.Float32()*2 - 1
	}

	s.IsShaking = true
}

func (s *ShakeStruct) Teardown() {
	s.t = 0
	s.IsShaking = false
}
func (s *ShakeStruct) Update(dt float32) {
	if !s.IsShaking {
		return
	}

	s.t += dt

	if s.t > s.duration {
		s.Teardown()
	}
}

func (s *ShakeStruct) Amplitude() float32 {
	sample := s.t * s.frequency
	s0 := int(math.Floor(float64(sample)))
	s1 := s0 + 1

	k := s.Decay()

	// low := s.noise(s1)
	// high := s.noise(s0)

	// fmt.Printf("s0: %v, s1: %v, low: %v, high: %v\n", s0, s1, low, high)
	m := s.noise(s1) - s.noise(s0)

	x := sample - float32(s0)
	b := s.noise(s0)
	y := m*x + b

	// result := s.noise(s0) + (sample-float32(s0))*(s.noise(s1)-s.noise(s0))*k
	// fmt.Println(result, "   ", y, "  ", x, "   ", m, "'   ", b, "   ", y)
	return y * (1 - k) * AMPLITUDE
}
func (s *ShakeStruct) noise(index int) float32 {
	if index >= len(s.samples) {
		// fmt.Println("WTFFF, ", len(s.samples))
		return 0
	}
	return s.samples[index]
}
func (s *ShakeStruct) Decay() float32 {
	return s.t / s.duration
}
