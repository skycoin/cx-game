package agent_draw

// import (
// 	"fmt"

// 	"github.com/go-gl/gl/v4.1-core/gl"
// 	"github.com/go-gl/mathgl/mgl32"

// 	indexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/IndexBufferDY"
// 	vertexbuffer "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/VertexBufferDY"
// 	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/geometry"
// 	vertexArray "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexArrayDY"
// 	vertexbufferLayout "github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/vertexbufferLayoutDY"
// 	"github.com/skycoin/cx-game/components/agents"
// 	"github.com/skycoin/cx-game/constants"
// 	"github.com/skycoin/cx-game/cxmath"
// 	"github.com/skycoin/cx-game/cxmath/math32"
// 	"github.com/skycoin/cx-game/engine/input"
// 	"github.com/skycoin/cx-game/engine/spine"
// 	"github.com/skycoin/cx-game/engine/spriteloader"
// 	"github.com/skycoin/cx-game/engine/spriteloader/spineloader"
// 	"github.com/skycoin/cx-game/physics/timer"
// 	"github.com/skycoin/cx-game/render"
// )

// func SpinePlayerDrawHandler3(agents []*agents.Agent, ctx DrawHandlerContext) {
// 	spineloader.Program.Use()
// 	defer spineloader.Program.StopUsing()

// 	//gl.Enable(gl.DEPTH_TEST)
// 	gl.Enable(gl.BLEND)
// 	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

// 	spineloader.Program.SetVec4F("u_Color", 0.8, 0.3, 0.8, 1.0)
// 	var offset uint32 = 0
// 	for i := 0; i < MaxIndexCount; i += 6 {
// 		indices[i+0] = 0 + offset
// 		indices[i+1] = 1 + offset
// 		indices[i+2] = 2 + offset

// 		indices[i+3] = 2 + offset
// 		indices[i+4] = 3 + offset
// 		indices[i+5] = 0 + offset

// 		offset += 4
// 	}

// 	//setup vertex array
// 	va = vertexArray.SetUpVertxArray()
// 	// setup and run vertex buffer
// 	vb = vertexbuffer.RunVertexBuffer(5 * 4 * MaxVertexCount)
// 	//setup vertex layout
// 	vbl = &vertexbufferLayout.VertexbufferLayout{}
// 	//add vertex buffer to vertex bufferlayout

// 	vbl.Push(gl.FLOAT, 2)
// 	vbl.Push(gl.FLOAT, 2)
// 	vbl.Push(gl.FLOAT, 1)
// 	va.AddBuffer(vb, vbl)

// 	// setup and run index buffer
// 	ib = indexbuffer.RunIndexBuffer(indices, len(indices))

// 	spineloader.Program.SetUniform1iv("u_Texture", int32(len(samplers)), &samplers[0])
// 	// var positions []float32

// 	for _, agent := range agents {

// 		body := &agent.Transform

// 		// bodyMultiplier := body.Size.Y / agent.Skeleton.Data.Size.Y

// 		alpha := timer.GetTimeBetweenTicks() / constants.MS_PER_TICK
// 		var interpolatedPos cxmath.Vec2
// 		if !agent.Transform.PrevPos.Equal(agent.Transform.Pos) {
// 			interpolatedPos = agent.Transform.PrevPos.Mult(1 - alpha).Add(agent.Transform.Pos.Mult(alpha))

// 		} else {
// 			interpolatedPos = agent.Transform.Pos
// 		}

// 		if body.Direction != 1 {
// 			agent.Skeleton.FlipX = false
// 		} else {
// 			agent.Skeleton.FlipX = true
// 		}

// 		if body.IsOnGround() && input.GetAxis(input.HORIZONTAL) != 0 {
// 			if moving == false {
// 				agent.SetAnimation(6)
// 				agent.Play = true
// 				moving = true
// 			}

// 		} else {
// 			if moving == true {
// 				agent.SetAnimation(2)
// 				moving = false
// 			}

// 		}

// 		//transform := translate.Mul4(scale)
// 		// agent.Skeleton.Local.Translate.Set(float32(body.Pos.X), float32(body.Pos.Y-(body.Size.Y/2)))
// 		//	agent.Skeleton.Local.Scale.Set(scaleX*0.025, agent.Transform.Size.Y*0.025)
// 		// agent.Skeleton.Local.Scale.Set(-1*bodyMultiplier, 1*bodyMultiplier)
// 		// agent.Skeleton.Local.Scale.Set(1, 1)

// 		// fmt.Println(agent.Skeleton.Data.Size)
// 		// fmt.Println(agent.Transform.Size)

// 		agent.Update(float64(alpha), cxmath.Vec2{body.Pos.X, body.Pos.Y})
// 		// agent.Skeleton.Local.Scale.Set(scaleX, agent.Transform.Size.Y)
// 		//agent.Transform.Size = cxmath.Vec2{agent.Skeleton.Data.Size.X, agent.Skeleton.Data.Size.Y}
// 		var pos []float32

// 		for i, slot := range agent.Skeleton.Order {

// 			bone := slot.Bone
// 			switch attachment := slot.Attachment.(type) {
// 			case nil:
// 			case *spine.RegionAttachment:
// 				if attachment.Name == "gun" {
// 					//	panic(fmt.Sprintf("Fpund %v", attachment))

// 					local := attachment.Local.Affine()
// 					final := bone.World.Mul(local)

// 					xform := geometry.Matrix(final.Col64())
// 					fmt.Println(xform)
// 					// spriteID := agent.Meta.GetImage(attachment.Name)
// 					spriteID := agent.Meta.GetImage(attachment.Name)
// 					m := render.GetSpriteByID(spriteID)
// 					//	spineloader.Program.Use()
// 					//gl.AC
// 					gl.BindTexture(gl.TEXTURE_2D, m.Texture.Texture)

// 					x, y := final.Translation().XY()

// 					// S_T := GetSpriteTransform(final)
// 					q := CreateQuad(float32(x), float32(y), float32(m.Height), float32(m.Width), float32(i), xform)
// 					pos = append(pos, q...)

// 					break
// 				}
// 				// spriteTransform := mgl32.Ident4().
// 				// 	Mul4(mgl32.Translate3D(S_T.Translation.X(), S_T.Translation.Y(), zOffset+float32(i))).
// 				// 	Mul4(mgl32.HomogRotate3DZ(S_T.Rotation)).

// 				// 	// Mul4(mgl32.Scale3D(-spriteSize[0]*0.065*body.Direction*hitboxToRender, spriteSize[1]*0.065, 1))
// 				// 	//Mul4(mgl32.Scale3D(S_T.Scale.X()*9, S_T.Scale.Y()*9, 1))
// 				// 	Mul4(mgl32.Scale3D(1/S_T.Scale.X()*bodyMultiplier*0.5, 1/S_T.Scale.Y()*bodyMultiplier*0.5, 1)).
// 				// 	Mul4(mgl32.ShearZ3D(S_T.Skew.X(), S_T.Skew.Y()))
// 				// set blending mode
// 				// BUG: incorrect, should use blending mode not compositing mode
// 				switch slot.Data.Blend {
// 				case spine.Normal:
// 					// MISSING
// 				case spine.Additive:
// 				//	draw.CompositeMode = ebiten.CompositeModeLighter
// 				case spine.Multiply:
// 					// MISSING
// 				case spine.Screen:
// 					// MISSING
// 					//	draw.CompositeMode = ebiten.CompositeModeLighter
// 				}
// 				// gl.BindVertexArray(render.QuadVao)

// 				// projection := spriteloader.Window.GetProjectionMatrix()
// 				// mvp := projection.Mul4(ctx.Camera.GetViewMatrix()).Mul4(wrappedTransform)

// 				// anim.Program.SetMat4("mvp", &mvp)
// 				// // texTransform := mgl32.Mat3
// 				// texTransform := agent.AnimationPlayback.Frame().Transform
// 				// anim.Program.SetMat3("texTransform", &texTransform)
// 				// gl.DrawArrays(gl.TRIANGLES, 0, 6)
// 				//_, _, _, _, tx, ty := draw.GeoM.GetElements()
// 				// poject := Project(mgl32.Vec2{final.Translation().X, final.Translation().Y}, final)

// 				// if attachment.Name == "goggles" {
// 				// 	fmt.Println("goggles sizex:", m.M_width)
// 				// 	fmt.Println("goggles sizey:", m.M_height)
// 				// 	fmt.Println("goggles transform:", final.Translation().Y)

// 				// }
// 				//if attachment.Name == "gun" || attachment.Name == "head" || attachment.Name == "goggles" {
// 				// q := CreateQuad(float32(m.M_height), float32(m.M_width), float32(i), xform, m)
// 				// pos = append(pos, q...)
// 				//}
// 				// fmt.Println(i)
// 				//offset += 50.0
// 				//target.DrawImage(m, &draw)

// 			default:
// 				panic(fmt.Sprintf("unknown attachment %v", attachment))
// 			}

// 		}

// 		vb.Bind()

// 		vb.BufferSubData(pos)

// 		translate := mgl32.Translate3D(
// 			interpolatedPos.X,
// 			interpolatedPos.Y,
// 			constants.AGENT_Z,
// 		)

// 		scale := mgl32.Scale3D(
// 			1,
// 			1,
// 			1,
// 		)
// 		transform := translate.Mul4(scale)
// 		wrappedTransform := wrapTransform2(
// 			transform,
// 			ctx.Camera.PlanetWidth,
// 			ctx.Camera.GetTransform(),
// 		)

// 		projection := spriteloader.Window.GetProjectionMatrix()
// 		mvp := projection.Mul4(ctx.Camera.GetViewMatrix()).Mul4(wrappedTransform)

// 		spineloader.Program.SetMat4("u_MVP", &mvp)
// 		va.Bind()
// 		ib.Bind()
// 		gl.DrawElements(gl.TRIANGLES, int32(ib.GetCount()), gl.UNSIGNED_INT, nil)
// 		// texTransform := agent.AnimationPlayback.Frame().Transform
// 		// anim.Program.SetMat3("texTransform", &texTransform)
// 		// gl.DrawArrays(gl.TRIANGLES, 0, 6)
// 	}
// }

// var QuadPosition [4]mgl32.Vec4
