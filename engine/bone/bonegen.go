package bonegen

import (
	"fmt"

	"github.com/go-gl/gl/v2.1/gl"
	animjson "github.com/skycoin/cx-game/engine/spriteloader/anim/json"
)

var (
	triangle = []float32{
		0.45, 0.46, 0, // top
		0, 0, 0, // left
		0.5, 0.5, 0, // right
	}
)

type LineBone struct {
	Name     string
	Parent   string
	Rotation float32
	Length   float32
	X        float32
	Y        float32
}

func DrawBone(vao uint32) {
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))
}

func GenerateBones(bones []animjson.Bone) {
	for _, lineBone := range bones {
		fmt.Println("bone: ", lineBone.Name)
		// DrawBone()
	}

}
