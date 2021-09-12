package agent_ai

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/cxmath/math32"
	"github.com/skycoin/cx-game/engine/spriteloader/anim"
	"github.com/skycoin/cx-game/physics"
	"github.com/skycoin/cx-game/world"
)

const (
	drillSpeed     float32 = 0.5
	drillJumpSpeed float32 = 15
)

func spiderDrilling(directionX float32, spiderDrill_PhysicsState *physics.Body, spiderDrill_playback *anim.Playback, planet *world.Planet) {
	spiderPos := spiderDrill_PhysicsState.Pos
	if directionX < 0 {
		spiderPos.X = spiderPos.X + directionX
	} else {
		spiderPos.X = spiderPos.X + (directionX * 2)
	}

	if planet.TileIsSolid(int(spiderPos.X), int(spiderPos.Y)) {
		tile, _ := planet.DamageTile(int(spiderPos.X), int(spiderPos.Y), world.TopLayer)
		tileHead, _ := planet.DamageTile(int(spiderPos.X), int(spiderPos.Y+1), world.TopLayer)
		_ = tile
		_ = tileHead
	}
}

func spiderAttack(spiderDrill_PhysicsState *physics.Body, spiderDrill_playback *anim.Playback) {
	spiderDrill_playback.PlayRepeating("Attack")
	spiderDrill_PhysicsState.Vel.X = 0
}

func AiHandlerDrill(agent *agents.Agent, ctx AiContext) {
	dist := ctx.PlayerPos.X() - agent.PhysicsState.Pos.X
	directionX := math32.Sign(dist)
	if math32.Abs(dist) > ctx.WorldWidth/2 {
		directionX *= -1
	}
	agent.PhysicsState.Direction = directionX * -1
	if math32.Abs(dist) > 1 {
		agent.AnimationPlayback.PlayRepeating("Walk")
		agent.PhysicsState.Vel.X = directionX * walkSpeed
	} else {
		// line of sigth
		isLoS := (ctx.PlayerPos.Y() - 0.5) == agent.PhysicsState.Pos.Y
		if isLoS {
			spiderAttack(&agent.PhysicsState, &agent.AnimationPlayback)
		}
	}

	isCollisionHorizontal := agent.PhysicsState.Collisions.Horizontal()
	if isCollisionHorizontal {
		// events.OnSpiderCollisionHorizontal.Trigger(events.SpiderEventData{
		// 	Agent: agent,
		// })
		spiderDrilling(directionX, &agent.PhysicsState, &agent.AnimationPlayback, &ctx.World.Planet)

	}

	// doJump :=
	// 	agent.PhysicsState.Collisions.Horizontal() &&
	// 		agent.PhysicsState.IsOnGround() && !agent.PhysicsState.Collisions.VerticalAbove()
	// if doJump {
	// 	events.OnSpiderBeforeJump.Trigger(events.SpiderEventData{
	// 		Agent: agent,
	// 	})

	// 	agent.PhysicsState.Vel.Y = drillJumpSpeed
	// 	// trigger an event when spiderdrill jump
	// 	events.OnSpiderJump.Trigger(events.SpiderEventData{
	// 		Agent: agent,
	// 	})
	// } else {
	// 	agent.PhysicsState.Vel.Y -= constants.Gravity * constants.PHYSICS_TICK
	// }
}
