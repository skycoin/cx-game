package agent_draw

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/spineMath"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/render"
)

func drawSpineSprite(
	agent *agents.Agent, ctx DrawHandlerContext,
	spriteID render.SpriteID, zOffset float32,
) {
	body := &agent.Transform

	//drawn one frame behind
	alpha := timer.GetTimeBetweenTicks() / constants.PHYSICS_TICK

	var interpolatedPos cxmath.Vec2
	if !body.Pos.Equal(body.PrevPos) {
		interpolatedPos = body.PrevPos.Mult(1 - alpha).Add(body.Pos.Mult(alpha))
	} else {
		interpolatedPos = body.Pos

	}
	// interpolatedPos.X = math32.Round(interpolatedPos.X*32) / 32
	translate := mgl32.Translate3D(
		interpolatedPos.X,
		interpolatedPos.Y,
		zOffset,
	)

	hitboxToRender := 1 / constants.PLAYER_RENDER_TO_HITBOX
	scaleX := -body.Size.X * body.Direction * hitboxToRender
	scale := mgl32.Scale3D(scaleX, agent.Transform.Size.Y, 2)

	transform := translate.Mul4(scale)
	render.DrawWorldSpineSprite(transform, spriteID, render.NewSpriteDrawOptions())

	for _, slot := range agent.Skeleton.Order {
		//win.SetComposeMethod(pixel.ComposeOver)
		/*
			switch slot.Data.Blend {
			case spine.Normal:
				win.SetComposeMethod(pixel.ComposeOver)
			case spine.Additive:
				// MISSING
				win.SetComposeMethod(pixel.ComposePlu)
			case spine.Multiply:
				// MISSING
				win.SetComposeMethod(pixel.ComposeOver)
			case spine.Screen:
				// MISSING
				win.SetComposeMethod(pixel.ComposePlus)
			default:
				win.SetComposeMethod(pixel.ComposeOver)
			}
		*/

		bone := slot.Bone
		switch attachment := slot.Attachment.(type) {
		case nil:
		case *spine.RegionAttachment:
			// final := bone.World.Mul(attachment.Local.Affine())
			// xform := pixel.Matrix(final.Col64())
			local := attachment.Local.Affine()
			final := bone.World.Mul(local)

			var geom spineMath.GeoM
			geom.SetElement(0, 0, float64(final.M00))
			geom.SetElement(0, 1, float64(final.M01))
			geom.SetElement(0, 2, float64(final.M02))
			geom.SetElement(1, 0, float64(final.M10))
			geom.SetElement(1, 1, float64(final.M11))
			geom.SetElement(1, 2, float64(final.M12))

			// scale2 := mgl32.Scale3D(attachment.Size.X, attachment.Size.Y, 1)
			// transform2 := translate.Mul4(scale2)

			//	pd := char.GetImage(attachment.Name, attachment.Path)
			//	sprite := pixel.NewSprite(pd, pd.Rect)
			switch {
			case attachment.Name == "back_arm":
				fmt.Println("test name: ", attachment.Name)

				fmt.Println("test size: ", attachment.Size)
				render.DrawWorldSprite(transform, agent.Meta.Back_armSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.Back_armSpriteID, constants.PLAYER_Z)

			case attachment.Name == "back_foot":
				render.DrawWorldSprite(transform, agent.Meta.Back_footSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.Back_footSpriteID, constants.PLAYER_Z)

			case attachment.Name == "back_hand":
				render.DrawWorldSprite(transform, agent.Meta.Back_handSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.Back_handSpriteID, constants.PLAYER_Z)

			case attachment.Name == "back_leg":
				render.DrawWorldSprite(transform, agent.Meta.Back_legSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.Back_legSpriteID, constants.PLAYER_Z)

			case attachment.Name == "body":
				render.DrawWorldSprite(transform, agent.Meta.BodySpriteID, render.NewSpriteDrawOptions())
				//drawSpineSprite(agent, ctx, agent.Meta.BodySpriteID, constants.PLAYER_Z)

			case attachment.Name == "front_arm":
				render.DrawWorldSprite(transform, agent.Meta.Front_armSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.Front_armSpriteID, constants.PLAYER_Z)

			case attachment.Name == "front_foot":
				render.DrawWorldSprite(transform, agent.Meta.Back_footSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.Front_footSpriteID, constants.PLAYER_Z)

			case attachment.Name == "front_hand":
				render.DrawWorldSprite(transform, agent.Meta.Front_handSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.Front_handSpriteID, constants.PLAYER_Z)

			case attachment.Name == "front_leg":
				render.DrawWorldSprite(transform, agent.Meta.Front_legtSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.Front_legtSpriteID, constants.PLAYER_Z)

			case attachment.Name == "head":

				render.DrawWorldSprite(transform, agent.Meta.HeadSpriteID, render.NewSpriteDrawOptions())
			//	drawSpineSprite(agent, ctx, agent.Meta.HeadSpriteID, constants.PLAYER_Z)

			default:
				//	panic("could not find character part")
			}

			//	sprite.DrawColorMask(win, xform, slot.Color)

		//	fmt.Println("test here: ", attachment)

		case *spine.MeshAttachment:
			fmt.Println("hit meshAttachment at: ", attachment)
			// pd := char.GetImage(attachment.Name, attachment.Path)
			// size := pd.Bounds().Size()

			// worldPosition := attachment.CalculateWorldVertices(char.Skeleton, slot)
			// tridata := pixel.MakeTrianglesData(len(attachment.Triangles) * 3)
			// for base, tri := range attachment.Triangles {
			// 	for k, index := range tri {
			// 		tri := &(*tridata)[base*3+k]
			// 		tri.Position = worldPosition[index].V()
			// 		uv := attachment.UV[index]
			// 		uv.Y = 1 - uv.Y
			// 		tri.Picture = uv.V()
			// 		tri.Picture = tri.Picture.ScaledXY(size)
			// 	}
			// }

			// var col pixel.RGBA
			// col.R, col.G, col.B, col.A = slot.Color.RGBA64()
			// for i := range *tridata {
			// 	tri := &(*tridata)[i]
			// 	tri.Color = col
			// 	tri.Intensity = 1
			// }

			// batch := pixel.NewBatch(tridata, pd)
			// batch.Draw(win)
		default:
			panic(fmt.Sprintf("unknown attachment %v", attachment))
		}
	}
	// render.DrawWorldSprite(transform, spriteID, render.NewSpriteDrawOptions())
}

// func (char *Agent) Draw(agents *Agent, ctx agent_draw.DrawHandlerContext) {
// 	// BUG: no way to get the current compose method for restoring
// 	//defer win.SetComposeMethod(pixel.ComposeOver)

// 	for _, slot := range char.Skeleton.Order {
// 		//win.SetComposeMethod(pixel.ComposeOver)
// 		/*
// 			switch slot.Data.Blend {
// 			case spine.Normal:
// 				win.SetComposeMethod(pixel.ComposeOver)
// 			case spine.Additive:
// 				// MISSING
// 				win.SetComposeMethod(pixel.ComposePlu)
// 			case spine.Multiply:
// 				// MISSING
// 				win.SetComposeMethod(pixel.ComposeOver)
// 			case spine.Screen:
// 				// MISSING
// 				win.SetComposeMethod(pixel.ComposePlus)
// 			default:
// 				win.SetComposeMethod(pixel.ComposeOver)
// 			}
// 		*/

// 		bone := slot.Bone
// 		switch attachment := slot.Attachment.(type) {
// 		case nil:
// 		case *spine.RegionAttachment:
// 			// final := bone.World.Mul(attachment.Local.Affine())
// 			// xform := pixel.Matrix(final.Col64())
// 			local := attachment.Local.Affine()
// 			final := bone.World.Mul(local)

// 			var geom spineMath.GeoM
// 			geom.SetElement(0, 0, float64(final.M00))
// 			geom.SetElement(0, 1, float64(final.M01))
// 			geom.SetElement(0, 2, float64(final.M02))
// 			geom.SetElement(1, 0, float64(final.M10))
// 			geom.SetElement(1, 1, float64(final.M11))
// 			geom.SetElement(1, 2, float64(final.M12))

// 			//	pd := char.GetImage(attachment.Name, attachment.Path)
// 			//	sprite := pixel.NewSprite(pd, pd.Rect)
// 			switch {
// 			case attachment.Name == "back_arm":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.Back_armSpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "back_foot":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.Back_footSpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "back_hand":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.Back_handSpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "back_leg":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.Back_legSpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "body":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.BodySpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "front_arm":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.Front_armSpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "front_foot":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.Front_footSpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "front_hand":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.Front_handSpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "front_leg":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.Front_legtSpriteID, constants.PLAYER_Z)

// 			case attachment.Name == "head":
// 				agent_draw.drawSpineSprite(agents, ctx, agents.Meta.HeadSpriteID, constants.PLAYER_Z)

// 			default:
// 				panic("could not find character part")
// 			}

// 			//	sprite.DrawColorMask(win, xform, slot.Color)

// 			fmt.Println("test here: ", attachment)

// 		case *spine.MeshAttachment:
// 			fmt.Println("hit meshAttachment at: ", attachment)
// 			// pd := char.GetImage(attachment.Name, attachment.Path)
// 			// size := pd.Bounds().Size()

// 			// worldPosition := attachment.CalculateWorldVertices(char.Skeleton, slot)
// 			// tridata := pixel.MakeTrianglesData(len(attachment.Triangles) * 3)
// 			// for base, tri := range attachment.Triangles {
// 			// 	for k, index := range tri {
// 			// 		tri := &(*tridata)[base*3+k]
// 			// 		tri.Position = worldPosition[index].V()
// 			// 		uv := attachment.UV[index]
// 			// 		uv.Y = 1 - uv.Y
// 			// 		tri.Picture = uv.V()
// 			// 		tri.Picture = tri.Picture.ScaledXY(size)
// 			// 	}
// 			// }

// 			// var col pixel.RGBA
// 			// col.R, col.G, col.B, col.A = slot.Color.RGBA64()
// 			// for i := range *tridata {
// 			// 	tri := &(*tridata)[i]
// 			// 	tri.Color = col
// 			// 	tri.Intensity = 1
// 			// }

// 			// batch := pixel.NewBatch(tridata, pd)
// 			// batch.Draw(win)
// 		default:
// 			panic(fmt.Sprintf("unknown attachment %v", attachment))
// 		}
// 	}

// 	// imd := imdraw.New(nil)
// 	// defer imd.Draw(win)

// 	// if char.DebugBones {
// 	// 	for _, bone := range char.Skeleton.Bones {
// 	// 		h := float64(bone.Data.Length)
// 	// 		if h < 10 {
// 	// 			h = 10
// 	// 		}

// 	// 		imd.SetMatrix(pixel.Matrix(bone.World.Col64()))
// 	// 		imd.Color = bone.Data.Color.WithAlpha(0.5)

// 	// 		w := h * 0.1
// 	// 		imd.Push(pixel.V(h, 0))
// 	// 		imd.Push(pixel.V(w+w, -w))
// 	// 		imd.Push(pixel.V(w+w, w))
// 	// 		imd.Polygon(0)

// 	// 		if bone.Parent != nil {
// 	// 			imd.SetMatrix(pixel.IM)
// 	// 			a := pixel.Vec(bone.World.Translation().V())
// 	// 			b := pixel.Vec(bone.Parent.World.Transform(spine.Vector{bone.Parent.Data.Length, 0}).V())
// 	// 			imd.Push(a)
// 	// 			imd.Push(b)
// 	// 			imd.Line(1)
// 	// 		}
// 	// 	}
// 	// }

// 	// if char.DebugCenter {
// 	// 	imd.SetMatrix(pixel.Matrix(char.Skeleton.World().Col64()))
// 	// 	imd.Color = pixel.RGB(0, 0, 1)
// 	// 	imd.Push(pixel.Vec{})
// 	// 	imd.Circle(20, 2)
// 	// }

// }

func SpineDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	for _, agent := range agents {

		//agent.Draw(agent, ctx)
		//drawPlayerSprite(agent, ctx, agent.Meta.SuitSpriteID, constants.PLAYER_Z)
		//	drawPlayerSprite(agent, ctx, agent.Meta.HelmetSpriteID, constants.PLAYER_Z+1)

		drawSpineSprite(agent, ctx, agent.Meta.Back_armSpriteID, constants.PLAYER_Z)
		//drawSpineSprite(agent, ctx, agent.Meta.Back_footSpriteID, constants.PLAYER_Z+1)
		//drawSpineSprite(agent, ctx, agent.Meta.Back_handSpriteID, constants.PLAYER_Z+2)
		//drawSpineSprite(agent, ctx, agent.Meta.Back_legSpriteID, constants.PLAYER_Z+3)
		//drawSpineSprite(agent, ctx, agent.Meta.BodySpriteID, constants.PLAYER_Z)
		//drawSpineSprite(agent, ctx, agent.Meta.Front_armSpriteID, constants.PLAYER_Z+5)
		//drawSpineSprite(agent, ctx, agent.Meta.Front_footSpriteID, constants.PLAYER_Z+6)
		//drawSpineSprite(agent, ctx, agent.Meta.Front_handSpriteID, constants.PLAYER_Z+7)
		//drawSpineSprite(agent, ctx, agent.Meta.Front_legtSpriteID, constants.PLAYER_Z+8)
		//drawSpineSprite(agent, ctx, agent.Meta.HeadSpriteID, constants.PLAYER_Z+1)
	}
}

//***************** spine test code **********************/
//***************************************/
