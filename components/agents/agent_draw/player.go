package agent_draw

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/engine/spriteloader"
)

const (
	playerHeadSize float32 = 1.5
)

func PlayerDrawHandler(agents []*agents.Agent, ctx DrawHandlerContext) {
	if len(agents) == 0 {
		return
	}
	drawOpts := spriteloader.NewDrawOptions()
	for _, agent := range agents {
		scaleX :=
			-agent.PhysicsState.Size.X * agent.PhysicsState.Direction

		spriteloader.DrawSpriteQuadOptions(
			agent.PhysicsState.Pos.X-ctx.Camera.X,
			agent.PhysicsState.Pos.Y-ctx.Camera.Y,
			scaleX, agent.PhysicsState.Size.Y,
			agent.PlayerData.SuitSpriteID, drawOpts,
		)
		spriteloader.DrawSpriteQuadOptions(
			agent.PhysicsState.Pos.X-ctx.Camera.X,
			agent.PhysicsState.Pos.Y-ctx.Camera.Y,
			scaleX, agent.PhysicsState.Size.Y,
			agent.PlayerData.HelmetSpriteID, drawOpts,
		)
	}
}
