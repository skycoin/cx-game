package particle_draw

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants/particle_constants"
)

func DrawTransparent(particleList []*particles.Particle, cam *camera.Camera) {
	shader := GetShader(particle_constants.DRAW_HANDLER_ALPHA_BLENDED)

	shader.Use()
	//accomplished by setting blendFunc
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	for _, particle := range particleList {
		world := mgl32.Translate3D(
			particle.Position.X,
			particle.Position.Y,
			0,
		)
		shader.SetMat4("world", &world)
		shader.SetVec4("color", &mgl32.Vec4{1, 1, 1,
			(particle.Duration - particle.TimeToLive) / particle.Duration,
		})
		shader.SetInt("particle_texture", 0)

		gl.BindTexture(gl.TEXTURE_2D, particle.Texture)
		gl.DrawArrays(gl.TRIANGLES, 0, 6)

	}

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
