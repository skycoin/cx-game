package agent_draw

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/render"
)

func spineDrawPlayerSprite(
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
	interpolatedPos.X = math32.Round(interpolatedPos.X*32) / 32

	translate := mgl32.Translate3D(
		interpolatedPos.X,
		interpolatedPos.Y,
		zOffset,
	)

	hitboxToRender := 1 / constants.PLAYER_RENDER_TO_HITBOX
	scaleX := -body.Size.X * body.Direction * hitboxToRender
	scale := mgl32.Scale3D(scaleX, agent.Transform.Size.Y, 1)

	agent.Skeleton.Data.Size.X = scaleX
	agent.Skeleton.Data.Size.Y = agent.Transform.Size.Y

	transform := translate.Mul4(scale)

	for _, slot := range agent.Skeleton.Order {
		bone := slot.Bone
		switch attachment := slot.Attachment.(type) {
		case nil:
		case *spine.RegionAttachment:
			local := attachment.Local.Affine()
			final := bone.World.Mul(local)

			spriteTranslate := mgl32.Translate3D(final.Translation().X, final.Translation().Y, zOffset)

			spriteScale := mgl32.Scale3D(1, 1, 1)

			spriteTransform := spriteTranslate.Mul4(spriteScale)

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

			render.DrawWorldSprite(spriteTransform, agent.Meta.GetImage(attachment.Name), render.NewSpriteDrawOptions())
			_, _, _, _, tx, ty := draw.GeoM.GetElements()
			poject := Project(mgl32.Vec2{final.Translation().X, final.Translation().Y}, final)

			if attachment.Name == "goggles" {
				fmt.Println("goggles sizex:", m.M_width)
				fmt.Println("goggles sizey:", m.M_height)
				fmt.Println("goggles transform:", final.Translation().Y)

			}
			if attachment.Name == "gun" || attachment.Name == "head" || attachment.Name == "goggles" {
				q := CreateQuad(float32(m.M_height), float32(m.M_width), float32(i), xform, m)
				pos = append(pos, q...)
			}
			fmt.Println(i)
			offset += 50.0
			target.DrawImage(m, &draw)
		default:
			panic(fmt.Sprintf("unknown attachment %v", attachment))
		}
	}

	render.DrawWorldSprite(transform, spriteID, render.NewSpriteDrawOptions())
}

func SpinePlayerDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	for _, agent := range agents {
		spineDrawPlayerSprite(agent, ctx, agent.Meta.SuitSpriteID, constants.PLAYER_Z)
		//	drawPlayerSprite(agent, ctx, agent.Meta.HelmetSpriteID, constants.PLAYER_Z+1)
	}
}
