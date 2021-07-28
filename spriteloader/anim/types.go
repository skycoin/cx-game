package anim

import (
	"github.com/go-gl/mathgl/mgl32"
)

type AnimationID int

type Animation struct {
	ID AnimationID
	Texture uint32
	Actions map[string]Action
}

type Action struct {
	Name string
	Frames []Frame
}

type Frame struct {
	Duration float32
	Transform mgl32.Mat3 // 2x2 transform needs 3d matrix for translation
}

type Playback struct {
	Animation Animation
	Tex uint32
	SecondsIntoFrame float32
	FrameIndex int

	Action,NextAction Action
	Repeat bool
}
