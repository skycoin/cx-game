package agent_draw

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/spriteloader"
)

func AnimatedDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	// TODO
	drawOpts := spriteloader.NewDrawOptions()
	for _, agent := range agents {
		spriteID := agent.AnimationState.Action.SpriteID()
		spriteloader.DrawSpriteQuadOptions(
			agent.PhysicsState.Pos.X-ctx.Camera.X,
			agent.PhysicsState.Pos.Y-ctx.Camera.Y,
			agent.PhysicsState.Size.X*agent.PhysicsState.Direction,
			agent.PhysicsState.Size.Y,
			spriteID, drawOpts,
		)
	}
}
