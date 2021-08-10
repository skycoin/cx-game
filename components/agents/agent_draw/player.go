package agent_draw

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/render"
)

const hitboxToRender float32 = 64.0/43

func drawPlayerSprite(
		agent *agents.Agent, ctx DrawHandlerContext,
		spriteID render.SpriteID, zOffset float32,
) {
	body := &agent.PhysicsState

	translate := mgl32.Translate3D( body.Pos.X, body.Pos.Y, zOffset )

	scaleX := -body.Size.X * body.Direction * hitboxToRender
	scale := mgl32.Scale3D( scaleX, agent.PhysicsState.Size.Y, 1)

	transform := translate.Mul4(scale)
	render.DrawWorldSprite(transform, spriteID, render.NewSpriteDrawOptions() )
}

func PlayerDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	for _, agent := range agents {
		drawPlayerSprite(agent, ctx, agent.PlayerData.SuitSpriteID, 0)
		drawPlayerSprite(agent, ctx, agent.PlayerData.HelmetSpriteID, 1)
	}
}
