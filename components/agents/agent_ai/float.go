package agent_ai

import (
	"math/rand"

	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/cxmath/math32"
)

const dashSpeed    float32 = 4
const dashCooldown float32 = 2
const maxHeight float32 = 20
const aggroRadius float32 = 20

func AiHandlerEnemyFloating(agent *agents.Agent, ctx AiContext) {
	if !agent.IsWaiting() {
		dash(agent,ctx)
		agent.Timers.WaitingFor = dashCooldown
	}
}

func dashAngle(agent *agents.Agent, ctx AiContext) float32 {
	displacement := ctx.PlayerPos.Sub(agent.Transform.Pos.Mgl32())
	isAggro := displacement.Len() < aggroRadius
	if isAggro {
		// aim at player +/- 45 degrees
		angleToPlayer := math32.Atan2(displacement.Y(),displacement.X())
		return angleToPlayer+math32.Pi/2*(rand.Float32()-0.5)
	}

	heightAt :=
		ctx.World.Planet.GetHeight(int(math32.Round(agent.Transform.Pos.X)))

	isTooHigh := ( agent.Transform.Pos.Y - float32(heightAt) ) > maxHeight
	// force Y to be downwards
	if (isTooHigh) { return math32.Pi+rand.Float32()*math32.Pi }

	return rand.Float32()*math32.Pi*2
}

func dash(agent *agents.Agent, ctx AiContext) {
	angle := dashAngle(agent, ctx)
	dashVelocity := cxmath.Vec2 {
		dashSpeed*math32.Sin(angle),
		dashSpeed*math32.Cos(angle),
	}
	agent.Transform.Vel = agent.Transform.Vel.Add(dashVelocity)
}
