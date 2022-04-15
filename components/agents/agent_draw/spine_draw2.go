package agent_draw

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/Texture"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/geometry"
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/input"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/engine/spriteloader"
	"github.com/skycoin/cx-game/engine/spriteloader/spineloader"
	"github.com/skycoin/cx-game/physics/timer"
)

var Z_OffSet float32 = 0.0
var samplers = [31]int32{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
}

func SpinePlayerDrawHandler2(agents []*agents.Agent, ctx DrawHandlerContext) {
	//temporary prevent from writing to depth buffer
	gl.DepthMask(false)

	gl.Disable(gl.DEPTH_TEST)
	//accomplished by setting blendFunc
	gl.Enable(gl.BLEND)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	spineloader.Shader.Bind()
	defer spineloader.Shader.UnBind()
	// renderer.GLClearError()
	// gl.Enable(gl.DEPTH_TEST)

	spineloader.Shader.SetUniForm4f("u_Color", 1.0, 1.0, 1.0, 1.0)

	spineloader.Shader.SetUniform1iv("u_Texture", int32(len(samplers)), &samplers[0])
	var positions []float32

	if input.GetKeyDown(glfw.KeyPageUp) {
		Z_OffSet = Z_OffSet + 10
		fmt.Println(Z_OffSet)
	}
	if input.GetKeyDown(glfw.KeyPageDown) {
		Z_OffSet = Z_OffSet - 10
		fmt.Println(Z_OffSet)
	}

	for _, agent := range agents {

		body := &agent.Transform

		bodyMultiplier := body.Size.Y / agent.Meta.Skeleton.Data.Size.Y

		alpha := timer.GetTimeBetweenTicks() / constants.MS_PER_TICK
		var interpolatedPos cxmath.Vec2
		if !body.Pos.Equal(body.PrevPos) {
			interpolatedPos = body.PrevPos.Mult(1 - alpha).Add(body.Pos.Mult(alpha))
		} else {
			interpolatedPos = body.Pos

		}

		if body.Direction != 1 {
			agent.Meta.Skeleton.FlipX = false
		} else {
			agent.Meta.Skeleton.FlipX = true
		}
		if body.IsOnGround() && input.GetAxis(input.HORIZONTAL) != 0 {
			if moving == false {
				agent.Meta.SetAnimation(3)
				agent.Meta.Play = true
				moving = true
			}

		} else {
			agent.Meta.SetAnimation(2)
			moving = false

		}

		//transform := translate.Mul4(scale)
		// agent.Meta.Skeleton.Local.Translate.Set(float32(body.Pos.X), float32(body.Pos.Y-(body.Size.Y/2)))
		//	agent.Skeleton.Local.Scale.Set(scaleX*0.025, agent.Transform.Size.Y*0.025)
		// agent.Meta.Skeleton.Local.Scale.Set(-1*bodyMultiplier, 1*bodyMultiplier)
		// agent.Skeleton.Local.Scale.Set(-1*0.077, 1*0.077)

		skeletonTraslate := cxmath.Vec2{0, float32(0)}
		skeletonScale := cxmath.Vec2{1 * bodyMultiplier, 1 * bodyMultiplier}
		// fmt.Println(agent.Skeleton.Data.Size)
		// fmt.Println(skeletonTraslate)

		agent.Meta.Update(float64(alpha), skeletonTraslate, skeletonScale)
		// agent.Skeleton.Local.Scale.Set(scaleX, agent.Transform.Size.Y)
		//agent.Transform.Size = cxmath.Vec2{agent.Skeleton.Data.Size.X, agent.Skeleton.Data.Size.Y}
		// var pos []float32

		for i, slot := range agent.Meta.Skeleton.Order {

			bone := slot.Bone
			switch attachment := slot.Attachment.(type) {
			case nil:
			case *spine.RegionAttachment:
				if attachment.Name != "" {
					//	panic(fmt.Sprintf("Fpund %v", attachment))

					local := attachment.Local.Affine()
					final := bone.World.Mul(local)

					m := agent.Meta.GetImage2(attachment.Name, attachment.Path)
					xform := geometry.Matrix(final.Col64())
					//fmt.Println(xform)
					// spriteID := agent.Meta.GetImage(attachment.Name)

					//m := render.GetSpriteByID(spriteID)
					//	spineloader.Program.Use()
					//gl.AC
					// gl.BindTexture(gl.TEXTURE_2D, m.Texture.Texture)

					x, y := final.Translation().XY()

					// spriteTransform := mgl32.Ident4().
					// 	Mul4(mgl32.Translate3D(S_T.Translation.X(), S_T.Translation.Y(), zOffset+float32(i))).
					// 	Mul4(mgl32.HomogRotate3DZ(S_T.Rotation)).

					// 	// Mul4(mgl32.Scale3D(-spriteSize[0]*0.065*body.Direction*hitboxToRender, spriteSize[1]*0.065, 1))
					// 	//Mul4(mgl32.Scale3D(S_T.Scale.X()*9, S_T.Scale.Y()*9, 1))
					// 	Mul4(mgl32.Scale3D(1/S_T.Scale.X()*bodyMultiplier*0.5, 1/S_T.Scale.Y()*bodyMultiplier*0.5, 1)).
					// 	Mul4(mgl32.ShearZ3D(S_T.Skew.X(), S_T.Skew.Y()))
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
					m.Bind(uint32(i))

					// S_T := GetSpriteTransform(final)
					q := CreateQuad(float32(x), float32(y), float32(m.M_height), float32(m.M_width), float32(i), xform, m)
					positions = append(positions, q...)
					// gl.BindVertexArray(render.QuadVao)

					// projection := spriteloader.Window.GetProjectionMatrix()
					// mvp := projection.Mul4(ctx.Camera.GetViewMatrix()).Mul4(wrappedTransform)

					// anim.Program.SetMat4("mvp", &mvp)
					// // texTransform := mgl32.Mat3
					// texTransform := agent.AnimationPlayback.Frame().Transform
					// anim.Program.SetMat3("texTransform", &texTransform)
					// gl.DrawArrays(gl.TRIANGLES, 0, 6)
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
				}
			default:
				panic(fmt.Sprintf("unknown attachment %v", attachment))
			}

		}

		spineloader.Ib.Bind()
		spineloader.Va.Bind()

		spineloader.Vb.Bind()

		spineloader.Vb.BufferSubData(positions)
		//	fmt.Println(interpolatedPos.X, interpolatedPos.Y)
		translate := mgl32.Translate3D(
			interpolatedPos.X,
			interpolatedPos.Y-(body.Size.Y/2),
			0,
		)
		scale := mgl32.Scale3D(
			1,
			1,
			1,
		)

		transform := translate.Mul4(scale)
		wrappedTransform := wrapTransform2(
			transform,
			ctx.Camera.PlanetWidth,
			ctx.Camera.GetTransform(),
		)

		projection := spriteloader.Window.GetProjectionMatrix()
		mvp := projection.Mul4(ctx.Camera.GetViewMatrix()).Mul4(wrappedTransform)

		//	spineloader.Shader.Bind()
		spineloader.Shader.SetUniFormMat4f("u_MVP", mvp)

		spineloader.Render.DrawDY(spineloader.Va, spineloader.Ib, spineloader.Shader)
		//spineloader.Vb.Unbind()
		gl.ActiveTexture(gl.TEXTURE0)

		/* // still under testing
		var posBone []float32
		if agent.Meta.DebugBones {
			for _, bone := range agent.Meta.Skeleton.Bones {
				h := float64(bone.Data.Length) * float64(bodyMultiplier)

				if h < (10 * float64(bodyMultiplier)) {
					h = 10 * float64(bodyMultiplier)
				}
				fmt.Println(h)
				triForm := geometry.Matrix(bone.World.Col64())
				color := bone.Data.Color.WithAlpha(1)
				spineloader.Shader.SetUniForm4f("u_Color", color.R, color.G, color.B, color.A)

				w := h * 0.7
				// imd.Push(geometry.V(h, 0))
				// imd.Push(geometry.V(w+w, -w))
				// imd.Push(geometry.V(w+w, w))
				// imd.Polygon(0)
				posBone = CreateTri(float32(h), float32(w), triForm)
				// if i == 2 {
				// 	break
				// }

				// if bone.Parent != nil {
				// 	imd.SetMatrix(geometry.IM)
				// 	a := geometry.Vec(bone.World.Translation().V())
				// 	b := geometry.Vec(bone.Parent.World.Transform(spine.Vector{bone.Parent.Data.Length, 0}).V())
				// 	imd.Push(a)
				// 	imd.Push(b)
				// 	imd.Line(1)
				// }
			}
		}

		spineloader.TriIb.Bind()
		spineloader.TriVa.Bind()

		spineloader.TriVb.Bind()

		spineloader.TriVb.BufferSubData(posBone)
		//	fmt.Println(interpolatedPos.X, interpolatedPos.Y)
		translateBones := mgl32.Translate3D(
			interpolatedPos.X,
			interpolatedPos.Y-(body.Size.Y/2),
			0,
		)
		scaleBones := mgl32.Scale3D(
			1,
			1,
			1,
		)

		transformBones := translateBones.Mul4(scaleBones)
		wrappedTransformBones := wrapTransform2(
			transformBones,
			ctx.Camera.PlanetWidth,
			ctx.Camera.GetTransform(),
		)

		projectionBones := spriteloader.Window.GetProjectionMatrix()
		mvpBones := projectionBones.Mul4(ctx.Camera.GetViewMatrix()).Mul4(wrappedTransformBones)
		fmt.Println(transformBones)
		//spineloader.Shader.Bind()
		spineloader.Shader.SetUniFormMat4f("u_MVP", mvpBones)

		spineloader.Render.DrawDY(spineloader.TriVa, spineloader.TriIb, spineloader.Shader)
		spineloader.TriVb.Unbind()
		spineloader.TriIb.Unbind()
		spineloader.TriVa.Unbind()
		*/

	}
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.DepthMask(true)
	//gl.Enable(gl.BLEND)
	// renderer.GLCheckError()
}

type Vertex struct {
	Position  mgl32.Vec2
	TexCoords mgl32.Vec2
	TexID     float32
}

// var QuadPosition [4]mgl32.Vec4

// func Project(u mgl32.Vec2, m spine.Affine) mgl32.Vec2 {
// 	return mgl32.Vec2{m.M00*u.X() + m.M01*u.Y() + m.M02, m.M10*u.X() + m.M11*u.Y() + m.M12}
// }

func CreateQuad(x, y, h, w, textureID float32, matrix geometry.Matrix, m *Texture.Texture) []float32 {

	var (
		horizontal = geometry.V(float64(w/2), 0)
		vertical   = geometry.V(0, float64(h/2))
	)

	var pos []float32
	v0 := Vertex{}
	v0.TexCoords = mgl32.Vec2{0.0, 0.0}
	v0.TexID = textureID
	t0 := geometry.Vec{}.Sub(horizontal).Sub(vertical)
	xy0 := matrix.Project(t0)
	v0.Position = mgl32.Vec2{float32(xy0.X), float32(xy0.Y)}
	pos = append(pos, v0.Position.X(), v0.Position.Y(), v0.TexCoords.X(), v0.TexCoords.Y(), v0.TexID)

	v1 := Vertex{}
	v1.TexCoords = mgl32.Vec2{1.0, 0.0}
	v1.TexID = textureID
	t1 := geometry.Vec{}.Add(horizontal).Sub(vertical)
	xy1 := matrix.Project(t1)
	v1.Position = mgl32.Vec2{float32(xy1.X), float32(xy1.Y)}
	pos = append(pos, v1.Position.X(), v1.Position.Y(), v1.TexCoords.X(), v1.TexCoords.Y(), v1.TexID)

	v2 := Vertex{}
	v2.TexCoords = mgl32.Vec2{1.0, 1.0}
	v2.TexID = textureID
	t2 := geometry.Vec{}.Add(horizontal).Add(vertical)
	xy2 := matrix.Project(t2)
	v2.Position = mgl32.Vec2{float32(xy2.X), float32(xy2.Y)}
	pos = append(pos, v2.Position.X(), v2.Position.Y(), v2.TexCoords.X(), v2.TexCoords.Y(), v2.TexID)

	v3 := Vertex{}
	v3.TexCoords = mgl32.Vec2{0.0, 1.0}
	v3.TexID = textureID
	t3 := geometry.Vec{}.Sub(horizontal).Add(vertical)
	xy3 := matrix.Project(t3)
	v3.Position = mgl32.Vec2{float32(xy3.X), float32(xy3.Y)}
	pos = append(pos, v3.Position.X(), v3.Position.Y(), v3.TexCoords.X(), v3.TexCoords.Y(), v3.TexID)

	return pos

}
func CreateTri(h, w float32, matrix geometry.Matrix) []float32 {

	var (
	// horizontal = geometry.V(float64(h), 0)
	// vertical   = geometry.V(0, float64(w))
	)

	var pos []float32
	v0 := Vertex{}
	v0.TexCoords = mgl32.Vec2{0.0, 0.0}
	v0.TexID = 0
	//	t0 := geometry.Vec{}.Sub(horizontal).Sub(vertical)
	t0 := geometry.Vec{float64(h), 0}
	xy0 := matrix.Project(t0)
	v0.Position = mgl32.Vec2{float32(t0.X), float32(xy0.Y)}
	pos = append(pos, v0.Position.X(), v0.Position.Y(), v0.TexCoords.X(), v0.TexCoords.Y(), v0.TexID)

	v1 := Vertex{}
	v1.TexCoords = mgl32.Vec2{1.0, 0.0}
	v1.TexID = 0
	// t1 := geometry.Vec{}.Add(horizontal).Sub(vertical)
	t1 := geometry.Vec{float64(w + w), float64(-w)}
	xy1 := matrix.Project(t1)
	v1.Position = mgl32.Vec2{float32(xy1.X), float32(xy1.Y)}
	pos = append(pos, v1.Position.X(), v1.Position.Y(), v1.TexCoords.X(), v1.TexCoords.Y(), v1.TexID)

	v2 := Vertex{}
	v2.TexCoords = mgl32.Vec2{1.0, 1.0}
	v2.TexID = 0
	// t2 := geometry.Vec{}.Add(horizontal).Add(vertical)
	t2 := geometry.Vec{float64(w + w), float64(w)}
	xy2 := matrix.Project(t2)
	v2.Position = mgl32.Vec2{float32(xy2.X), float32(xy2.Y)}
	pos = append(pos, v2.Position.X(), v2.Position.Y(), v2.TexCoords.X(), v2.TexCoords.Y(), v2.TexID)

	return pos

}

func wrapTransform2(raw mgl32.Mat4, worldWidth float32, cameraTransform mgl32.Mat4) mgl32.Mat4 {
	rawX := raw.At(0, 3)
	x := math32.PositiveModulo(rawX, worldWidth)
	camX := cameraTransform.At(0, 3)
	if x-camX > worldWidth/2 {
		x -= worldWidth
	}
	if x-camX < -worldWidth/2 {
		x += worldWidth
	}

	translate := mgl32.Translate3D(x-rawX, 0, 0)
	return translate.Mul4(raw)
}
