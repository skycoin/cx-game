package starmap

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"

	//"gopkg.in/yaml.v2"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/render"
	sl "github.com/skycoin/cx-game/spriteloader"
)

var texture uint32
var gradient uint32
var program uint32
var vao uint32
var window *render.Window

var xpos float32 = 0.0
var ypos float32 = 0.0
var xwidth float32 = 6
var yheight float32 = 6

func Init(_window *render.Window) {
	window = _window
}

func Generate(size int, scale float32, levels uint8) {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	perlinField := perlin.NewPerlin2D(1, 512, 4, 256)
	max := float32(math.Sqrt2 / 1.9)
	min := float32(-math.Sqrt2 / 1.9)

	// Set color for each pixel.
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			val := perlinField.Noise(float32(x)*scale, float32(y)*scale, 0.5, 2, 8)
			val = (val - min) / (max - min)                           // normalized aproximation
			brightness := uint8(val*float32(levels)) * (255 / levels) // map values
			//brightness := uint8(val * 255)
			img.Set(x, y, color.RGBA{brightness, brightness, brightness, 255})
		}
	}

	//file, _ := os.Create("test_noise.png")
	//png.Encode(file, img)

	fragSource, err := ioutil.ReadFile("./starmap/gradient.glsl")
	if err != nil {
		panic(err)
	}
	vertSource, err := ioutil.ReadFile("./starmap/vertex.glsl")
	if err != nil {
		panic(err)
	}

	fragment, err := render.CompileShader(string(fragSource), gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	vertex, err := render.CompileShader(string(vertSource), gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	texture = sl.MakeTexture(img)
	_, gradientImg := sl.LoadPng("./assets/starfield/gradients/heightmap_gradient_07.png")
	gradient = sl.MakeTexture(gradientImg)

	program = gl.CreateProgram()

	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)
	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	vao = sl.MakeQuadVao()
}

func Draw() {
	gl.UseProgram(program) // use shader

	textureLocation := gl.GetUniformLocation(program, gl.Str("nebulaText\x00"))
	if textureLocation == -1 {
		log.Println("[WARNING] can't get the nebule texture uniform location")
	}
	gradientLocation := gl.GetUniformLocation(program, gl.Str("gradientText\x00"))
	if gradientLocation == -1 {
		log.Println("[WARNING] can't get the gradient texture uniform location")
	}

	gl.Uniform1i(textureLocation, 0)
	gl.Uniform1i(gradientLocation, 1)

	gl.Uniform2f(
		gl.GetUniformLocation(program, gl.Str("texScale\x00")),
		1.0, 1.0,
	)
	gl.Uniform2f(
		gl.GetUniformLocation(program, gl.Str("texOffset\x00")),
		float32(0.0), float32(0.0),
	)

	worldTranslate := mgl32.Mat4.Mul4(
		mgl32.Translate3D(float32(xpos), float32(ypos), -10),
		mgl32.Scale3D(float32(xwidth), float32(yheight), 1),
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(program, gl.Str("world\x00")),
		1, false, &worldTranslate[0],
	)

	aspect := float32(window.Width) / float32(window.Height)
	projectTransform := mgl32.Perspective(
		mgl32.DegToRad(45), aspect, 0.1, 100.0,
	)
	gl.UniformMatrix4fv(
		gl.GetUniformLocation(program, gl.Str("projection\x00")),
		1, false, &projectTransform[0],
	)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.ActiveTexture(gl.TEXTURE0 + 1)
	gl.BindTexture(gl.TEXTURE_2D, gradient)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
