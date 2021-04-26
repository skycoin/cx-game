package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime"

	"github.com/skycoin/cx-game/world"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	// "github.com/go-gl/mathgl/f32"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/render"
)

const (
	width  = 800
	height = 480
)

var (
	square = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,

		1, -1, 0, 1, 1,
		-1, -1, 0, 0, 1,
		-1, 1, 0, 0, 0,
	}
)
var tex uint32
var wx, wy, wz float32
var numOfStars int

func main() {
	runtime.LockOSThread()

	wx = 0
	wy = 0
	wz = -10
	win := render.NewWindow(height, width, true)
	window := win.Window
	defer glfw.Terminate()
	program := win.Program

	vao := makeVao(square)
	gl.GenTextures(2, &tex)
	stars := prepareStars()
	for !window.ShouldClose() {
		drawStarField(vao, window, program, stars)
	}
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, err := png.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	im := img.(*image.NRGBA)
	sim := im.SubImage(image.Rect(18, 9, 32, 22))
	// w := img.Bounds().Max.X
	// h := img.Bounds().Max.Y
	spriteH := 32
	spriteW := 32

	pixels := make([]byte, spriteW*spriteH*4)
	bIndex := 0
	for y := 0; y < spriteH; y++ {
		for x := 0; x < spriteW; x++ {
			r, g, b, a := sim.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(spriteW),
		int32(spriteH),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(pixels))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture, nil
}

func newTexture1(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, err := png.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	im := img.(*image.NRGBA)
	sim := im.SubImage(image.Rect(0, 9, 32, 32))
	// w := img.Bounds().Max.X
	// h := img.Bounds().Max.Y
	spriteH := 32
	spriteW := 32

	pixels := make([]byte, spriteW*spriteH*4)
	bIndex := 0
	for y := 0; y < spriteH; y++ {
		for x := 0; x < spriteW; x++ {
			r, g, b, a := sim.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(spriteW),
		int32(spriteH),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(pixels))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture, nil
}

func drawStarField(vao uint32, window *glfw.Window, program uint32, stars []world.Star) {
	// gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	// Load the texture
	texture, err := newTexture("./../../assets/starfield/stars/Starsheet1.png")
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < numOfStars; i++ {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i))
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture\x00")), 0)
		worldTranslate := mgl32.Translate3D(stars[i].X, stars[i].Y, wz)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("world\x00")), 1, false, &worldTranslate[0])
		projectTransform := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)
		gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projectTransform[0])
		gl.BindVertexArray(vao)
		gl.Enable(gl.BLEND);
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);		
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func prepareStars() []world.Star {
	stars := []world.Star{}

	numOfStars = rand.Intn(51) + 50 // todo get from commandline
	for i := 0; i < numOfStars; i++ {
		sx := rand.Float32() * 12 - 6 
		sy := rand.Float32() * 7 - 3
		star := world.Star{nil, sx, sy, -10}
		stars = append(stars, star)
	}

	return stars
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)

	return vao
}
