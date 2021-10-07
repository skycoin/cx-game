package agent_draw

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/physics/timer"
	"github.com/skycoin/cx-game/render"
)

func drawPlayerSprite(
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
	scale := mgl32.Scale3D(scaleX, agent.Transform.Size.Y, 1)

	transform := translate.Mul4(scale)
	render.DrawWorldSprite(transform, spriteID, render.NewSpriteDrawOptions())
}

func PlayerDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	for _, agent := range agents {
		drawPlayerSprite(agent, ctx, agent.Meta.SuitSpriteID, constants.PLAYER_Z)
		drawPlayerSprite(agent, ctx, agent.Meta.HelmetSpriteID, constants.PLAYER_Z+1)
	}
}
