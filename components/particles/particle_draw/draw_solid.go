package particle_draw

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants/particle_constants"
)

func DrawSolid(particleList []*particles.Particle, cam *camera.Camera) {
	shader := GetShader(particle_constants.DRAW_HANDLER_SOLID)
	shader.Use()
	projection := cam.GetProjectionMatrix()
	shader.SetMat4("projection", &projection)
	//todo, right now both draw handlers use same shader
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
}
