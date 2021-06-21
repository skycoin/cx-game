package particles

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/spriteloader"
)

type Bullet struct {
	transform mgl32.Mat4
	velocity mgl32.Vec2
}

var (
	bulletShader *utility.Shader
	bullets []Bullet
)

func InitBullets() {
	bulletShader = utility.NewShader(
		"./assets/shader/simple.vert", "./assets/shader/color.frag" )
	bulletShader.Use()
	bulletShader.SetVec4F("colour",1,0,0,1)
	bulletShader.StopUsing()
	bullets = []Bullet {}
}

func CreateBullet( origin mgl32.Vec2, velocity mgl32.Vec2 ) {
	log.Print("creating bullet")
	bullets = append(bullets, Bullet {
		transform: mgl32.Translate3D(origin.X(), origin.Y(), 0),
		velocity: velocity,
	})
}

func (bullet Bullet) draw(ctx render.Context) {
	log.Printf("drawing bullet at %v", bullet.transform.Col(3).Vec2())
	world := ctx.World.Mul4(bullet.transform)
	bulletShader.SetMat4("world", &world)

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

	gl.BindVertexArray(spriteloader.QuadVao);
}

func DrawBullets(ctx render.Context) {
	log.Print("drawing bullets")
	bulletShader.Use()
	bulletShader.SetMat4("projection", &ctx.Projection)
	configureGlForBullet()
	// TODO add texture
	// bulletShader.SetUint("tex", bulletShader.gpuTex)

	for _,bullet := range bullets {
		bullet.draw(ctx)
	}
	bulletShader.StopUsing()
}

//func TickBullets(ctx 
