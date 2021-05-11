package ui;
// runtime text rendering library.
// uses a single VAO/VBO pair
// and draws a different set of triangles depending on the character.

import (
	"log"
//	"math"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/spriteloader"
)

const fontTexWidth = 256
const fontTexHeight = 256
const fontScale = float32(10)

type NormalizedCharData struct {
	size mgl32.Vec2
	offset mgl32.Vec2
	index int
}

var asciiToCharDataMap = make(map[int]NormalizedCharData)
// opengl objects
var fontTex,vao,vbo uint32

// allocate a VBO for the entire font which can 
// render different characters with very little computation
func initFontVbo() {
	var vertexAttributes = make([]float32,5*6*len(charDatas))
	i := 0
	for _,charData := range charDatas {
		top := float32(charData.ty + charData.h)/256
		bottom := float32(charData.ty)/256
		right := float32(charData.tx + charData.w)/256
		left := float32(charData.tx)/256

		w := float32(charData.w)/256
		h := float32(charData.h)/256

		// tri 1
		vertexAttributes[i] = w
		vertexAttributes[i+1] = h
		vertexAttributes[i+2] = 0
		vertexAttributes[i+3] = right
		vertexAttributes[i+4] = bottom
		i += 5

		vertexAttributes[i] = w
		vertexAttributes[i+1] = 0
		vertexAttributes[i+2] = 0
		vertexAttributes[i+3] = right
		vertexAttributes[i+4] = top
		i += 5

		vertexAttributes[i] = 0
		vertexAttributes[i+1] = h
		vertexAttributes[i+2] = 0
		vertexAttributes[i+3] = left
		vertexAttributes[i+4] = bottom
		i += 5

		// tri 2
		vertexAttributes[i] = w
		vertexAttributes[i+1] = 0
		vertexAttributes[i+2] = 0
		vertexAttributes[i+3] = right
		vertexAttributes[i+4] = top
		i += 5

		vertexAttributes[i] = 0
		vertexAttributes[i+1] = 0
		vertexAttributes[i+2] = 0
		vertexAttributes[i+3] = left
		vertexAttributes[i+4] = top
		i += 5

		vertexAttributes[i] = 0
		vertexAttributes[i+1] = h
		vertexAttributes[i+2] = 0
		vertexAttributes[i+3] = left
		vertexAttributes[i+4] = bottom
		i += 5
	}
	log.Print(vertexAttributes[6*5:6*5+6*5])

	gl.GenBuffers(1,&vbo)
	gl.GenVertexArrays(1,&vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER,vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		4*len(vertexAttributes),
		gl.Ptr(vertexAttributes),
		gl.STATIC_DRAW,
	)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)
}

func InitTextRendering() {
	_,img,_ := spriteloader.LoadPng("assets/font/8bitoperator_jve.png")
	fontTex = spriteloader.MakeTexture(img)

	for _,charData := range charDatas {
		asciiToCharDataMap[charData.ascii] = NormalizedCharData {
			size: mgl32.Vec2 {
				float32(charData.w)/fontTexWidth,
				float32(charData.h)/fontTexHeight,
			},
			offset: mgl32.Vec2 {
				float32(charData.tx)/fontTexWidth,
				float32(charData.ty)/fontTexHeight,
			},
			index: charData.index,
		}
	}

	initFontVbo()
}

func calculateLineWidth(text string) float32 {
	x := float32(0)
	for _,charCode := range text {
		x += asciiToCharDataMap[int(charCode)].size.X()
	}
	return float32(x)*fontScale
}
// TODO line wrapping
// TODO alignment options
func DrawStringLeftAligned(text string, transform mgl32.Mat4) {
	// setup GPU params
	gl.Disable(gl.DEPTH_TEST)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, fontTex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// use the default program with support 
	// for scaled / offset texture lookups
	// TODO dedicate a program just for this
	program := spriteloader.Window.Program
	gl.UseProgram(program)
	// this line only needs to occur once if we have a dedicated program
	gl.Uniform1ui(
		gl.GetUniformLocation(program, gl.Str("ourTexture\x00")),
		fontTex,
	)

	// center it
	pos := mgl32.Vec2 {}
	for _, charCode := range text {
		charData,ok := asciiToCharDataMap[int(charCode)]
		if ok {
			letterTransform := transform.
				Mul4(mgl32.Translate3D(pos.X(),pos.Y(),0)).
				Mul4(cxmath.Scale(fontScale))

			gl.UniformMatrix4fv(
				gl.GetUniformLocation(program, gl.Str("world\x00")),
				1, false, &letterTransform[0],
			)
			aspect := float32(spriteloader.Window.Width) / float32(spriteloader.Window.Height)
			projectTransform := mgl32.Perspective(
				mgl32.DegToRad(45), aspect, 0.1, 100.0,
			)
			gl.UniformMatrix4fv(
				gl.GetUniformLocation(program, gl.Str("projection\x00")),
				1, false, &projectTransform[0],
			)
			gl.BindVertexArray(vao)
			glStart := 6*charData.index
			gl.DrawArrays(gl.TRIANGLES, int32(glStart), 6)
		}
		// TODO variable width fonts
		pos = pos.Add(mgl32.Vec2{charData.size.X()*fontScale,0})
	}
}

func DrawStringRightAligned(text string, transform mgl32.Mat4) {
	transformLeftAligned := transform.
		Mul4(mgl32.Translate3D(-calculateLineWidth(text),0,0))

	DrawStringLeftAligned(text, transformLeftAligned)
}
