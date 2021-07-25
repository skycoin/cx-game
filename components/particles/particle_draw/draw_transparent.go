package particle_draw

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
)

var quad_vao uint32

func DrawTransparent(particleList []*particles.Particle, cam *camera.Camera) {
	if quad_vao == 0 {
		vertices := []float32{
			-0.5, -0.5, 0, 1,
			-0.5, 0.5, 0, 0,
			0.5, -0.5, 1, 1,

			-0.5, 0.5, 0, 0,
			0.5, -0.5, 1, 1,
			0.5, 0.5, 1, 0,
		}
		gl.GenVertexArrays(1, &quad_vao)
		gl.BindVertexArray(quad_vao)

		var vao uint32
		gl.GenBuffers(1, &vao)
		gl.BindBuffer(gl.ARRAY_BUFFER, vao)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))

	}
	shader := GetShader(constants.PARTICLE_DRAW_HANDLER_TRANSPARENT)

	shader.Use()
	//accomplished by setting blendFunc
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	for _, particle := range particleList {
		world := mgl32.Translate3D(
			particle.Position.X-cam.X,
			particle.Position.Y-cam.Y,
			0,
		).Mul4(mgl32.Scale3D(1, 1, 1))
		// projection := mgl32.Ortho2D(0, 800.0/32, 0, 600.0/32)
		projection := mgl32.Ortho2D(
			-800.0/2/32, 800.0/2/32,
			-600.0/2/32, 600.0/2/32,
		)
		shader.SetMat4("projection", &projection)
		shader.SetMat4("world", &world)
		shader.SetVec4("color", &mgl32.Vec4{1, 1, 1,
			particle.TimeToLive / particle.Duration,
		})
		shader.SetInt("particle_texture", 0)

		gl.BindTexture(gl.TEXTURE_2D, particle.Texture)
		gl.BindVertexArray(quad_vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
