package particles

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
)

const bulletDamage = 1

type Bullet struct {
	transform mgl32.Mat4
	velocity  mgl32.Vec2
}

var (
	bulletProgram render.Program
	bullets       []Bullet
)

func InitBullets() {
	return 
	bulletProgram = render.CompileProgram(
		"./assets/shader/mvp.vert", "./assets/shader/color.frag")
	bulletProgram.Use()
	bulletProgram.SetVec4F("colour", 1, 0, 0, 1)
	bulletProgram.StopUsing()
	bullets = []Bullet{}
}

func CreateBullet(origin mgl32.Vec2, velocity mgl32.Vec2) {
	bullets = append(bullets, Bullet{
		transform: mgl32.Translate3D(origin.X(), origin.Y(), 0),
		velocity:  velocity,
	})

	//todo call emitter only for effect
	
}

func (bullet Bullet) WorldTransform() mgl32.Mat4 {
	return bullet.transform.Mul4(cxmath.Scale(1.0 / 4))
}

func (bullet Bullet) draw(ctx render.WorldContext) {
	mvp := ctx.ModelToModelViewProjection(bullet.WorldTransform())
	bulletProgram.SetMat4("mvp", &mvp)

	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func configureGlForBullet() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_2D)
	gl.ActiveTexture(gl.TEXTURE0)
	/*
		gl.BindTexture(gl.TEXTURE_2D, bulletTex.gpuTex)
		// blurry is better than jagged for a bullet
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	*/

	gl.BindVertexArray(render.QuadVao)
}

func DrawBullets(ctx render.WorldContext) {
	bulletProgram.Use()
	projection := ctx.Projection()
	bulletProgram.SetMat4("projection", &projection)
	configureGlForBullet()
	// TODO add texture
	// bulletProgram.SetUint("tex", bulletProgram.gpuTex)

	for _, bullet := range bullets {
		bullet.draw(ctx)
	}
	bulletProgram.StopUsing()
}

func TickBullets(dt float32) {
	// TODO replace with more robust projectile system
	// newBullets := []Bullet{}
	// for _, bullet := range bullets {
	// 	//bullet.ttl -= dt
	// 	bullet.transform = bullet.transform.Mul4(
	// 		mgl32.Translate3D(bullet.velocity.X()*dt, bullet.velocity.Y()*dt, 0))

	// 	collision, collided := physics.CheckCollision(bullet.WorldTransform())
	// 	_ = collision
	// 	if collided {
	// 		if collision.Body.Damage != nil {
	// 			collision.Body.Damage(bulletDamage)
	// 		}
	// 	} else {
	// 		newBullets = append(newBullets, bullet)
	// 	}
	// }
	// bullets = newBullets

}
