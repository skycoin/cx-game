package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	perlin "github.com/skycoin/cx-game/procgen"
	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
)

var (
	vertices    []float32
	myPerlin    perlin.Perlin2D
	noiseConfig *noiseSettings
)

func main() {

	myPerlin = perlin.NewPerlin2D(
		noiseConfig.seed,
		noiseConfig.x,
		noiseConfig.xs,
		noiseConfig.gradMax,
	)

	win := render.NewWindow(800, 600, false)
	window := win.Window
	shader := utility.NewShader("./test/perlin_timemap/shaders/vertex.glsl", "./test/perlin_timemap/shaders/fragment.glsl")
	shader.Use()
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}

func genVector(x, y int) []float32 {
	z := myPerlin.Noise()
}

func genData() []float32 {
	var data []float32
	for i := 0; i < 800; i++ {
		for j := 0; j < 600; j++ {
			data = append(data, genVector(i, j)...)
		}
	}

	return data
}

type noiseSettings struct {
	seed        int64
	x           int
	xs          int
	gradMax     int
	persistence float32
	lacunarity  float32
	octaves     int
}
