package agent_draw

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/geometry"
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/render"
)

var moving bool

func spineDrawPlayerSprite(
	agent *agents.Agent, ctx DrawHandlerContext,
	spriteID render.SpriteID, zOffset float32,
) {

	body := &agent.Transform
	agent.Skeleton.Data.Size.X = body.Size.X
	agent.Skeleton.Data.Size.Y = body.Size.Y

	//drawn one frame behind
	alpha := timer.GetTimeBetweenTicks() / constants.MS_PER_TICK

	var interpolatedPos cxmath.Vec2
	if !body.Pos.Equal(body.PrevPos) {
		interpolatedPos = body.PrevPos.Mult(1 - alpha).Add(body.Pos.Mult(alpha))
	} else {
		interpolatedPos = body.Pos

	}
	interpolatedPos.X = math32.Round(interpolatedPos.X*32) / 32

	// translate := mgl32.Translate3D(
	// 	interpolatedPos.X,
	// 	interpolatedPos.Y,
	// 	zOffset,
	// )

	hitboxToRender := 1 / constants.PLAYER_RENDER_TO_HITBOX
	// scaleX := -body.Size.X * body.Direction * hitboxToRender
	scaleX := -body.Size.X * hitboxToRender
	// scale := mgl32.Scale3D(scaleX, agent.Transform.Size.Y, 1)

	// agent.Skeleton.Data.Size.X = scaleX
	// agent.Skeleton.Data.Size.Y = agent.Transform.Size.Y

	if body.Direction != 1 {
		agent.Skeleton.FlipX = true
	} else {
		agent.Skeleton.FlipX = false
	}

	if body.IsOnGround() && input.GetAxis(input.HORIZONTAL) != 0 {
		if moving == false {
			agent.SetAnimation(3)
			agent.Play = true
			moving = true
		}

	} else {
		if moving == true {
			agent.SetAnimation(2)
			moving = false
		}

	}

	//transform := translate.Mul4(scale)
	agent.Skeleton.Local.Translate.Set(float32(interpolatedPos.X), float32(interpolatedPos.Y-1.5))
	//	agent.Skeleton.Local.Scale.Set(scaleX*0.025, agent.Transform.Size.Y*0.025)
	agent.Skeleton.Local.Scale.Set(scaleX*0.025, agent.Transform.Size.Y*0.025)

	agent.Update(float64(alpha), cxmath.Vec2{interpolatedPos.X, interpolatedPos.Y})
	// agent.Skeleton.Local.Scale.Set(scaleX, agent.Transform.Size.Y)
	agent.Transform.Size = cxmath.Vec2{agent.Skeleton.Data.Size.X, agent.Skeleton.Data.Size.Y}
	for i, slot := range agent.Skeleton.Order {
		bone := slot.Bone
		switch attachment := slot.Attachment.(type) {
		case nil:
		case *spine.RegionAttachment:
			spriteID := agent.Meta.GetImage(attachment.Name)

			// spriteSize := agent.Meta.GetImageSize(spriteID)
			// tex := render.GetSpriteByID(spriteID).Texture.Texture

			local := attachment.Local.Affine()
			final := bone.World.Mul(local)

			S_T := GetSpriteTransform(final)
			//xform := geometry.Matrix(final.Col64())

			/*
				x, y := final.Translation().XY()

				render.SpineQuadVao = render.MakeQuadVao(CreateQuad(float32(x), float32(y), spriteSize[1], spriteSize[0], xform))

				gl.BindVertexArray(render.SpineQuadVao)
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, tex)

				translate := mgl32.Translate3D(
					interpolatedPos.X,
					interpolatedPos.Y,
					constants.AGENT_Z+float32(i),
				)
				scale := mgl32.Scale3D(
					100,
					100,
					1,
				)
				transform := translate.Mul4(scale)
				wrappedTransform := wrapTransform(
					transform,
					ctx.Camera.PlanetWidth,
					ctx.Camera.GetTransform(),
				)
				projection := spriteloader.Window.GetProjectionMatrix()
				mvp := projection.Mul4(ctx.Camera.GetViewMatrix()).Mul4(wrappedTransform)

				anim.Program.SetMat4("mvp", &mvp)
				gl.DrawArrays(gl.TRIANGLES, 0, 6)
			*/
			// spriteTransform := mgl32.Ident4().
			// 	Mul4(mgl32.Translate3D(final.M00*final.M02+final.M01*final.M12,
			// 		final.M10*final.M02+final.M11*final.M12, constants.MIDLAYER_Z+1)).
			// 	Mul4(mgl32.HomogRotate3DZ(0)).
			// 	Mul4(mgl32.Scale3D(1, 1, 1))

			//fmt.Println(spriteSize)

			// translate := mgl32.Translate3D(
			// 	S_T.Translation.X(),
			// 	S_T.Translation.Y(),
			// 	zOffset+float32(i),
			// )

			// scale := mgl32.Scale3D(S_T.Scale.X()*12, -S_T.Scale.Y()*12, 1)

			// spriteTransform := translate.Mul4(scale)

			spriteTransform := mgl32.Ident4().
				Mul4(mgl32.Translate3D(S_T.Translation.X(), S_T.Translation.Y(), zOffset+float32(i))).
				Mul4(mgl32.HomogRotate3DZ(S_T.Rotation * body.Direction)).

				// Mul4(mgl32.Scale3D(-spriteSize[0]*0.065*body.Direction*hitboxToRender, spriteSize[1]*0.065, 1))
				Mul4(mgl32.Scale3D(S_T.Scale.X()*11, S_T.Scale.Y()*11, 1))

			//fmt.Println(spriteTransform)
			// spriteTransform := translate.Mul4()

			// spriteTranslate := mgl32.Translate3D(final.Translation().X, final.Translation().Y, zOffset)

			// spriteScale := mgl32.Scale3D(1, 1, 1)

			// spriteTransform := spriteTranslate.Mul4(spriteScale)

			//	xform := geometry.Matrix(final.Col64())
			//box := m.Bounds()

			// set blending mode
			// BUG: incorrect, should use blending mode not compositing mode
			switch slot.Data.Blend {
			case spine.Normal:
				// MISSING
			case spine.Additive:
			//	draw.CompositeMode = ebiten.CompositeModeLighter
			case spine.Multiply:
				// MISSING
			case spine.Screen:
				// MISSING
				//	draw.CompositeMode = ebiten.CompositeModeLighter
			}

			render.DrawWorldSprite(spriteTransform, spriteID, render.NewSpriteDrawOptions())
			//_, _, _, _, tx, ty := draw.GeoM.GetElements()
			// poject := Project(mgl32.Vec2{final.Translation().X, final.Translation().Y}, final)

			// if attachment.Name == "goggles" {
			// 	fmt.Println("goggles sizex:", m.M_width)
			// 	fmt.Println("goggles sizey:", m.M_height)
			// 	fmt.Println("goggles transform:", final.Translation().Y)

			// }
			//if attachment.Name == "gun" || attachment.Name == "head" || attachment.Name == "goggles" {
			// q := CreateQuad(float32(m.M_height), float32(m.M_width), float32(i), xform, m)
			// pos = append(pos, q...)
			//}
			// fmt.Println(i)
			//offset += 50.0
			//target.DrawImage(m, &draw)
		default:
			panic(fmt.Sprintf("unknown attachment %v", attachment))
		}
	}

	//	render.DrawWorldSprite(transform, spriteID, render.NewSpriteDrawOptions())
}

type Vertex struct {
	Position  mgl32.Vec2
	TexCoords mgl32.Vec2
}
type SpriteTransform struct {
	Translation mgl32.Vec2
	Rotation    float32
	Scale       mgl32.Vec2
	Skew        mgl32.Vec2
}

func GetSpriteTransform(mat spine.Affine) SpriteTransform {
	var a = mat.M00
	var b = mat.M01
	var c = mat.M10
	var d = mat.M11
	var e = mat.M02
	var f = mat.M12

	var delta = a*d - b*c

	result := SpriteTransform{
		Translation: mgl32.Vec2{e, f},
		Rotation:    0,
		Scale:       mgl32.Vec2{0, 0},
		Skew:        mgl32.Vec2{0, 0},
	}

	// Apply the QR-like decomposition.
	if a != 0 || b != 0 {

		var r = float32(math.Sqrt(float64(a*a + b*b)))

		// result.Rotation = b > 0 ? math.acos(a / r) : -math.acos(a / r);
		if b > 0 {
			result.Rotation = float32(math.Acos(float64(a / r)))
		} else {
			result.Rotation = float32(-math.Acos(float64(a / r)))
		}
		result.Scale = mgl32.Vec2{r, delta / r}
		result.Skew = mgl32.Vec2{float32(math.Atan(float64((a*c + b*d) / (r * r)))), 0}
	} else if c != 0 || d != 0 {
		var s = float32(math.Sqrt(float64(c*c + d*d)))
		//   result.Rotation =	math.Pi / 2 - (d > 0 ? math.acos(-c / s) : -math.acos(c / s));
		if d > 0 {
			result.Rotation = math.Pi/2 - float32(math.Acos(float64(-c/s)))
		} else {
			result.Rotation = math.Pi/2 - float32(-math.Acos(float64(-c/s)))
		}
		result.Scale = mgl32.Vec2{delta / s, s}
		result.Skew = mgl32.Vec2{0, float32(math.Atan(float64((a*c + b*d) / (s * s))))}
	} else {
		// a = b = c = d = 0
	}

	return result
}

func CreateQuad(x, y, h, w float32, matrix geometry.Matrix) []float32 {

	var (
		horizontal = geometry.V(float64(w/2), 0)
		vertical   = geometry.V(0, float64(h/2))
	)
	var pos []float32
	v0 := Vertex{}
	v0.TexCoords = mgl32.Vec2{0.0, 0.0}
	t0 := geometry.Vec{}.Sub(horizontal).Sub(vertical)
	xy0 := matrix.Project(t0)
	v0.Position = mgl32.Vec2{float32(xy0.X), float32(xy0.Y)}
	pos = append(pos, v0.Position.X(), v0.Position.Y(), 0, v0.TexCoords.X(), v0.TexCoords.Y())

	v1 := Vertex{}
	v1.TexCoords = mgl32.Vec2{1.0, 0.0}
	t1 := geometry.Vec{}.Add(horizontal).Sub(vertical)
	xy1 := matrix.Project(t1)
	v1.Position = mgl32.Vec2{float32(xy1.X), float32(xy1.Y)}
	pos = append(pos, v1.Position.X(), v1.Position.Y(), 0, v1.TexCoords.X(), v1.TexCoords.Y())

	v2 := Vertex{}
	v2.Position = mgl32.Vec2{x + w, y + h}
	v2.TexCoords = mgl32.Vec2{1.0, 1.0}
	t2 := geometry.Vec{}.Add(horizontal).Add(vertical)
	xy2 := matrix.Project(t2)
	v2.Position = mgl32.Vec2{float32(xy2.X), float32(xy2.Y)}
	pos = append(pos, v2.Position.X(), v2.Position.Y(), 0, v2.TexCoords.X(), v2.TexCoords.Y())

	v3 := Vertex{}
	v3.Position = mgl32.Vec2{x, y + h}
	v3.TexCoords = mgl32.Vec2{0.0, 1.0}
	t3 := geometry.Vec{}.Sub(horizontal).Add(vertical)
	xy3 := matrix.Project(t3)
	v3.Position = mgl32.Vec2{float32(xy3.X), float32(xy3.Y)}
	pos = append(pos, v3.Position.X(), v3.Position.Y(), 0, v3.TexCoords.X(), v3.TexCoords.Y())

	fmt.Printf("numbers=%v\n", pos)
	return pos

}

func SpinePlayerDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	// anim.Program.Use()
	// defer anim.Program.StopUsing()

	// gl.Enable(gl.DEPTH_TEST)

	for _, agent := range agents {
		spineDrawPlayerSprite(agent, ctx, agent.Meta.SuitSpriteID, constants.PLAYER_Z)
		//	drawPlayerSprite(agent, ctx, agent.Meta.HelmetSpriteID, constants.PLAYER_Z+1)
	}
}
