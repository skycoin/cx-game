package bonegen

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
)

type Bone struct {
	Name     string
	Parent   string
	Rotation float32
	Length   float32
	X        float32
	Y        float32
}

func DrawBone(x1, y1, x2, y2 float32) {
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(x1, y1)
	gl.Vertex2f(x2, y2)
	gl.End()
}

func GenerateBones(bones []Bone) {
	for _, bone := range bones {
		fmt.Println("bone: ", bone)
	}
}
