package particle_draw

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
)

func DrawSolid(particleList []*particles.Particle, cam *camera.Camera) {
	program := GetProgram(constants.PARTICLE_DRAW_HANDLER_SOLID)
	program.Use()
	defer program.StopUsing()

	gl.Disable(gl.BLEND)
	for _, particle := range particleList {
		world := mgl32.Translate3D(
			particle.Body.Pos.X-cam.X,
			particle.Body.Pos.Y-cam.Y,
			0,
		).Mul4(mgl32.Scale3D(1, 1, 1))
		// projection := mgl32.Ortho2D(0, 800.0/32, 0, 600.0/32)
		projection := mgl32.Ortho2D(
			-800.0/2/32, 800.0/2/32,
			-600.0/2/32, 600.0/2/32,
		)
		program.SetMat4("projection", &projection)
		program.SetMat4("world", &world)
		program.SetVec4("color", &mgl32.Vec4{1, 1, 1,
			particle.TimeToLive / particle.Duration,
		})
		program.SetInt("particle_texture", 0)

		gl.BindTexture(gl.TEXTURE_2D, particle.Texture)
		gl.BindVertexArray(quad_vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

	}
}
