package starmap

import (
	"image"
	"image/color"
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/render"
	sl "github.com/skycoin/cx-game/spriteloader"
	"gopkg.in/yaml.v2"
)

type noiseSettings struct {
	Size     int
	Scale    float32
	Levels   uint8
	Contrast float32

	Seed        int64
	Gradmax     int
	X           int
	Xs          int
	Persistance float32
	Lacunarity  float32
	Octaves     int

	GradFile string
}

type quadProp struct {
	xpos    float32
	ypos    float32
	xwidth  float32
	yheight float32
}

var texture uint32
var gradient uint32
var program uint32
var vao uint32
var window *render.Window

var frameCounter = 0 // reload the texture if yaml has change
const maxFrames = 60 // every 60 frames
const settingsFile = "./starmap/config/starmap.yaml"
const fragShaderFile = "./starmap/gradient.glsl"
const vertShaderFile = "./starmap/vertex.glsl"

var quad quadProp = quadProp{
	xpos:    0.0,
	ypos:    0.0,
	xwidth:  15.0,
	yheight: 10.0,
}

var settingsFileInfo os.FileInfo
var lastGradFile string

func Init(_window *render.Window) {
	window = _window
}

func Generate(size int, scale float32, levels uint8) {
	var err error
	settingsFileInfo, err = os.Stat(settingsFile)
	if err != nil {
		log.Panic(err)
	}
	loadTextures()

	fragSource, err := ioutil.ReadFile(fragShaderFile)
	if err != nil {
		log.Panic(err)
	}
	vertSource, err := ioutil.ReadFile(vertShaderFile)
	if err != nil {
		log.Panic(err)
	}

	fragment, err := render.CompileShader(string(fragSource), gl.FRAGMENT_SHADER)
	if err != nil {
		log.Panic(err)
	}
	vertex, err := render.CompileShader(string(vertSource), gl.VERTEX_SHADER)
	if err != nil {
		log.Panic(err)
	}

	program = gl.CreateProgram()

	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)
	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	vao = sl.MakeQuadVao()
}

func Draw() {
	frameCounter += 1
	if frameCounter > maxFrames {
		frameCounter = 0
		if fileHasChanged(settingsFile, settingsFileInfo) {
			var err error
			settingsFileInfo, err = os.Stat(settingsFile)
			if err != nil {
				log.Panic(err)
			}
			loadTextures()
		}
	}

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
		mgl32.Translate3D(float32(quad.xpos), float32(quad.ypos), -10),
		mgl32.Scale3D(float32(quad.xwidth), float32(quad.yheight), 1),
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

func loadTextures() {
	noise := noiseSettings{}

	settingsData, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		log.Panic(err)
	}
	err = yaml.Unmarshal(settingsData, &noise)
	if err != nil {
		log.Panic(err)
	}

	if lastGradFile != noise.GradFile {
		lastGradFile = noise.GradFile
		_, gradImg, _ := sl.LoadPng(noise.GradFile)
		gradient = sl.MakeTexture(gradImg)
	}

	img := genImage(noise)
	texture = sl.MakeTexture(img)
}

func genImage(noise noiseSettings) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, noise.Size, noise.Size))

	perlinField := perlin.NewPerlin2D(noise.Seed, noise.X, noise.Xs, noise.Gradmax)
	max := float32(math.Sqrt2 / (1.9 * noise.Contrast))
	min := float32(-math.Sqrt2 / (1.9 * noise.Contrast))

	// Set color for each pixel.
	for x := 0; x < noise.Size; x++ {
		for y := 0; y < noise.Size; y++ {
			val := perlinField.Noise(
				float32(x)*noise.Scale,
				float32(y)*noise.Scale,
				noise.Persistance,
				noise.Lacunarity,
				noise.Octaves,
			)
			val = clamp((val-min)/(max-min), 0.0, 1.0)                            // normalized aproximation
			brightness := uint8(val*float32(noise.Levels)) * (255 / noise.Levels) // map values
			//brightness := uint8(val * 255)
			img.Set(x, y, color.RGBA{brightness, brightness, brightness, 255})
		}
	}

	return img
}

func fileHasChanged(filepath string, fileInfo os.FileInfo) bool {
	stat, err := os.Stat(filepath)

	if err != nil {
		log.Panic(err)
	}

	if stat.Size() != fileInfo.Size() || stat.ModTime() != fileInfo.ModTime() {
		return true
	}

	return false
}

func clamp(number, min, max float32) float32 {
	if number > max {
		return max
	}
	if number < min {
		return min
	}
	return number
}
