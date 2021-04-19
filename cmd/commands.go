package commands

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/render"
	"github.com/urfave/cli/v2"
)

var tex uint32

func makeVao() uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 8*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(4*3))
	gl.EnableVertexAttribArray(1)

	return vao
}

var (
	vertices = []float32{
		1, 1, 0, 1, 0,
		1, -1, 0, 1, 1,
		-1, 1, 0, 0, 0,

		1, -1, 0, 1, 1,
		-1, -1, 0, 0, 1,
		-1, 1, 0, 0, 0,
	}
)

//commands
func Commands() {

	app := &cli.App{
		Name:  "stars",
		Usage: "Create a new window for the stars to show in",
		Action: func(c *cli.Context) error {
			err := glfw.Init()
			if err != nil {
				panic(err)
			}

			win := render.NewWindow(600, 800, true)
			window := win.Window
			defer glfw.Terminate()
			VAO := makeVao()
			program := win.Program
			gl.GenTextures(1, &tex)
			for !window.ShouldClose() {
				starsDraw(window, program, VAO)
			}

			// window, err := glfw.CreateWindow(640, 480, "Stars", nil, nil)
			// if err != nil {
			// 	panic(err)
			// }

			// window.MakeContextCurrent()

			// for !window.ShouldClose() {

			// 	window.SwapBuffers()
			// 	glfw.PollEvents()
			// }
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func starsDraw(window *glfw.Window, program uint32, VAO uint32) {

	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	imgFile, err := os.Open("./assets/cat.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer imgFile.Close()

	// Decode detexts the type of image as long as its image/<type> is imported
	im, err := png.Decode(imgFile)
	if err != nil {
		panic(err)
	}
	img := image.NewRGBA(im.Bounds())
	draw.Draw(img, img.Bounds(), im, image.Pt(0, 0), draw.Src)
	size := img.Rect.Size()
	gl.Enable(gl.TEXTURE_2D)
	gl.Enable(gl.BLEND)
	// gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(size.X), int32(size.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("ourTexture\x00")), 0)
	worldTranslate := mgl32.Translate3D(1, 1, 1)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("world\x00")), 1, false, &worldTranslate[0])
	projectTransform := mgl32.Perspective(mgl32.DegToRad(45), float32(1)/float32(1), 0.1, 100.0)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("projection\x00")), 1, false, &projectTransform[0])
	gl.BindVertexArray(VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
	glfw.PollEvents()
	window.SwapBuffers()
}
