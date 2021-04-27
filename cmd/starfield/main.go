package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/render"
	"github.com/urfave/cli/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Star struct {
	// Drawable uint32
	X     float32
	Y     float32
	Size  float32
	Depth float32
}

type Bitmap_RGBA struct {
	Width  int
	Height int
	Pixels []byte
}

const (
	width  = 800
	height = 600
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
	tex        uint32
	wx, wy, wz float64 = 0, 0, -10
	size       float64 = 1
	star       *Star
)

func main() {
	runtime.LockOSThread()

	initArgs()

	star = &Star{
		X:    float32(wx),
		Y:    float32(wy),
		Size: float32(size),
	}

	win := render.NewWindow(height, width, true)
	window := win.Window
	window.SetKeyCallback(keyCallback)
	defer glfw.Terminate()
	program := win.Program

	vao := makeVao(square)
	gl.GenTextures(1, &tex)
	for !window.ShouldClose() {
		gl.ClearColor(0, 0, 0.3, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)

		drawStarField(star, vao, window, program)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
func keyCallback(w *glfw.Window, k glfw.Key, scancode int, a glfw.Action, m glfw.ModifierKey) {
	if a != glfw.Press {
		return
	}
	if k == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
	switch k {
	case glfw.KeyUp:
		star.Y += 1
	case glfw.KeyDown:
		star.Y -= 1
	case glfw.KeyLeft:
		star.X -= 1
	case glfw.KeyRight:
		star.X += 1
	case glfw.KeyX:
		star.Size -= 0.1
	case glfw.KeyZ:
		star.Size += 0.1
	}
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	im, err := png.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	img := image.NewRGBA(im.Bounds())
	draw.Draw(img, img.Bounds(), im, image.Pt(0, 0), draw.Src)

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
		int32(img.Rect.Dx()),
		int32(img.Rect.Dy()),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(img.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture, nil
}

func drawStarField(star *Star, vao uint32, window *glfw.Window, program uint32) {

	// Load the texture
	texture, err := newTexture("./assets/starfield/stars/Starsheet1.png")

	if err != nil {
		log.Fatalln(err)
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture\x00")), 0)
	worldTranslate := mgl32.Translate3D(star.X, star.Y, float32(wz)/star.Size)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("world\x00")), 1, false, &worldTranslate[0])
	projectTransform := mgl32.Perspective(mgl32.DegToRad(45), float32(width)/float32(height), 0.1, 100.0)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projectTransform[0])
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))

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

func initArgs() {

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.Float64Flag{
			Name:        "x",
			Usage:       "x position for a star",
			Destination: &wx,
		},
		&cli.Float64Flag{
			Name:        "y",
			Usage:       "y position for a star",
			Destination: &wy,
		},
		&cli.Float64Flag{
			Name:        "size",
			Usage:       "size for a star (default = 1)",
			Value:       1,
			Destination: &size,
		},
	}
	app.Action = func(c *cli.Context) error {
		return nil
	}
	app.Run(os.Args)
	fmt.Printf("%f %f %f\n", wx, wy, wz)
}
