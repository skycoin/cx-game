package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/world"
	//cv "github.com/skycoin/cx-game/cmd/spritetool"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

var (
	DrawCollisionBoxes = false
	FPS                int
)

var CurrentPlanet *world.Planet

const (
	width  = 800
	height = 480
)

type Star struct {
	ImagePath                                      string
	X, Y                                           float32
	TexCoordX1, TexCoordY1, TexCoordX2, TexCoordY2 int
	Size                                           int32
	Depth                                          float32
}

func NewStar() Star {
	X := rand.Intn(3) * 8
	Y := rand.Intn(3) * 8

	star := Star{
		ImagePath:  "../../assets/starfield/stars/Starsheet1.png",
		X:          (rand.Float32() * 15) - (15 / 2),
		Y:          (rand.Float32() * 10) - (10 / 2),
		TexCoordX1: X,
		TexCoordY1: Y,
		TexCoordX2: X + 8,
		TexCoordY2: Y + 8,
		Size:       8,
	}

	return star
}

var (
	sprite = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,

		1, -1, 0, 1, 1,
		-1, -1, 0, 0, 1,
		-1, 1, 0, 0, 0,
	}
)

var wz float32
var Cam camera.Camera
var tex uint32

func makeVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(sprite), gl.Ptr(sprite), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)

	return vao
}

func main() {
	wz = -15
	win := render.NewWindow(height, width, true)
	window := win.Window
	defer glfw.Terminate()

	program := win.Program
	gl.GenTextures(1, &tex)

	var stars []Star
	for i := 0; i < 100; i++ {
		star := NewStar()
		stars = append(stars, star)
	}

	for !window.ShouldClose() {
		Tick()
		redraw(window, program, stars)
	}
}

func Tick() {

}

func redraw(window *glfw.Window, program uint32, stars []Star) {
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	for i := 0; i < len(stars); i++ {
		VAO := makeVao()
		drawStar(program, stars[i])
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func drawStar(program uint32, star Star) {
	imgFile, err := os.Open(star.ImagePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()

	// Decode detexts the type of image as long as its image/<type> is imported
	im, err := png.Decode(imgFile)
	if err != nil {
		log.Fatalln(err)
	}
	img := image.NewRGBA(image.Rect(star.TexCoordX1, star.TexCoordY1, star.TexCoordX2, star.TexCoordY2))
	draw.Draw(img, img.Bounds(), im, image.Pt(0, 0), draw.Src)
	size := img.Rect.Size()

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture\x00")), 0)
	worldTranslate := mgl32.Translate3D(star.X, star.Y, wz)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("world\x00")), 1, false, &worldTranslate[0])
	projectTransform := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projectTransform[0])
}
