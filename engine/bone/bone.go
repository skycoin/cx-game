package Bones

import { 
	"github.com/go-gl/gl/v4.1-core/gl"
}

type Bone struct {
	X1, Y1, X2, Y2 float32
	Father         *Bone
}

func (b *Bone) DrawBone(x1, y1, x2, y2 float32) {
	fmt.Println("DrawBone")
	gl.Begin(gl.LINE_STRIP)
	gl.Vertex2f(x1, y1)
	gl.Vertex2f(x2, y2)
	gl.End()
}
