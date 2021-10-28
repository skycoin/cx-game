package particle_draw

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/render"
)

var (
	quad_vao      uint32
	instanced_vao uint32
	data_vbo      uint32
)

const (
	MAX_PARTICLE_COUNT = 10000
)

func DrawTransparent(particleList []*particles.Particle, cam *camera.Camera) {
	program := GetProgram(constants.PARTICLE_DRAW_HANDLER_TRANSPARENT)

	program.Use()

	//temporary prevent from writing to depth buffer
	gl.DepthMask(false)

	//accomplished by setting blendFunc
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	for _, particle := range particleList {
		if particle == nil {
			continue
		}

		world := mgl32.Translate3D(
			particle.Pos.X,
			particle.Pos.Y,
			0,
		).Mul4(cxmath.Scale(particle.Size.X))
		projection := render.Projection
		program.SetMat4("projection", &projection)
		view := cam.GetViewMatrix()
		program.SetMat4("view", &view)
		program.SetMat4("world", &world)
		program.SetVec4("color", &mgl32.Vec4{1, 1, 1,
			particle.TimeToLive / particle.Duration,
		})
		program.SetInt("particle_texture", 0)

		gl.BindTexture(gl.TEXTURE_2D, particle.Texture)
		gl.BindVertexArray(quad_vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DepthMask(true)
}

func DrawTransparentInstanced(particleList []*particles.Particle, cam *camera.Camera) {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	//temporary prevent from writing to depth buffer
	gl.DepthMask(false)
	// gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE, gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ZERO)

	program := GetProgram(constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED)
	program.Use()

	projection := render.Projection
	program.SetMat4("projection", &projection)
	view := cam.GetViewMatrix()
	program.SetMat4("view", &view)
	program.SetInt("particle_texture", 0)

	bins := binByTexture(particleList)
	for texture, bin := range bins {
		gl.BindTexture(gl.TEXTURE_2D, texture)
		data_list := make([]float32, 0, len(bin)*4)

		for _, particle := range bin {
			data_list = append(data_list,
				particle.Pos.X, particle.Pos.Y, particle.Size.X,
				particle.TimeToLive/particle.Duration,
			)
		}

		updateBuffers(&data_list)

		gl.BindVertexArray(instanced_vao)
		gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(len(bin)))
	}

	gl.DepthMask(true)
}

func binByTexture(
	Particles []*particles.Particle,
) map[uint32][]*particles.Particle {
	bins := make(map[uint32][]*particles.Particle)
	for _, particle := range Particles {
		bins[particle.Texture] = append(bins[particle.Texture], particle)
	}
	return bins
}

func initDrawInstanced() {
	vertices := []float32{
		-0.5, -0.5, 0, 1,
		-0.5, 0.5, 0, 0,
		0.5, -0.5, 1, 1,

		-0.5, 0.5, 0, 0,
		0.5, -0.5, 1, 1,
		0.5, 0.5, 1, 0,
	}

	gl.GenVertexArrays(1, &instanced_vao)
	gl.BindVertexArray(instanced_vao)

	var quadvao uint32
	gl.GenBuffers(1, &quadvao)
	gl.BindBuffer(gl.ARRAY_BUFFER, quadvao)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(0, 0)

	gl.GenBuffers(1, &data_vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, data_vbo)
	gl.BufferData(gl.ARRAY_BUFFER, MAX_PARTICLE_COUNT*4, nil, gl.STREAM_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.VertexAttribDivisor(1, 1)

}

func updateBuffers(data_list *[]float32) {
	gl.BindVertexArray(instanced_vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, data_vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(*data_list)*4, gl.Ptr(*data_list), gl.STREAM_DRAW)
}
