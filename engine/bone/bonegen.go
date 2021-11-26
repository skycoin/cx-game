package bonegen

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
)

type Bone struct {
	X1, Y1, X2, Y2 float32
	Father         *Bone
}

type Bones []Bone

func DrawBone(x1, y1, x2, y2 float32) {
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(x1, y1)
	gl.Vertex2f(x2, y2)
	gl.End()
}

func GenerateBones(bones Bones) {
	for bone := range bones {
		fmt.Println("bone: ", bone)
	}
}
