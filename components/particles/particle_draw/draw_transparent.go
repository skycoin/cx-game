package particle_draw

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
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
	//accomplished by setting blendFunc
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	for _, particle := range particleList {
		world := mgl32.Translate3D(
			particle.Pos.X-cam.X,
			particle.Pos.Y-cam.Y,
			0,
		).Mul4(mgl32.Scale3D(1, 1, 1))
		// projection := mgl32.Ortho2D(0, 800.0/32, 0, 600.0/32)
		// projection := mgl32.Ortho2D(
		// 	-800.0/2/32, 800.0/2/32,
		// 	-600.0/2/32, 600.0/2/32,
		// )
		projection := cam.GetProjectionMatrix()
		program.SetMat4("projection", &projection)
		view := cam.GetView()
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
}

func DrawTransparentInstanced(particleList []*particles.Particle, cam *camera.Camera) {
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	program := GetProgram(constants.PARTICLE_DRAW_HANDLER_TRANSPARENT_INSTANCED)
	program.Use()

	projection := cam.GetProjectionMatrix()
	program.SetMat4("projection", &projection)
	program.SetInt("particle_texture", 0)

	gl.BindTexture(gl.TEXTURE_2D, particleList[0].Texture)
	data_list := make([]float32, 0, len(particleList)*4)

	for _, particle := range particleList {
		data_list = append(data_list, particle.Pos.X-cam.X, particle.Pos.Y-cam.Y, particle.Size.X, particle.TimeToLive/particle.Duration)
	}

	updateBuffers(&data_list)

	gl.BindVertexArray(instanced_vao)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(len(particleList)))
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
