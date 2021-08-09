package components

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/agents/agent_ai"
	"github.com/skycoin/cx-game/components/agents/agent_draw"
	"github.com/skycoin/cx-game/components/agents/agent_health"
	"github.com/skycoin/cx-game/components/agents/agent_physics"
	"github.com/skycoin/cx-game/components/particles/particle_draw"
	"github.com/skycoin/cx-game/components/particles/particle_physics"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/world"
)

func Update(dt float32) {
	updateTimers(currentWorld.Entities.Agents.Get(), dt)

	//update lifetimes
	currentWorld.Entities.Particles.Update(dt)

	emitter.SetPosition(currentPlayer.PhysicsState.Pos)
}

func updateTimers(agents []*agents.Agent, dt float32) {
	for _, agent := range agents {
		if agent.DrawHandlerID == constants.DRAW_HANDLER_ANIM {
			agent.AnimationPlayback.Update(dt)
		}
		if agent.WaitingFor > 0 {
			agent.WaitingFor -= dt
		}
		if agent.Died() {
			agent.TimeSinceDeath += dt
		}
	}
}

func FixedUpdate() {
	agent_health.UpdateAgents(&currentWorld.Entities.Agents)
	agent_physics.UpdateAgents(currentWorld)
	agent_ai.UpdateAgents(currentWorld, currentPlayer)
	particle_physics.Update(currentWorld)
}

func Draw(entities *world.Entities, cam *camera.Camera) {
	agent_draw.DrawAgents(&entities.Agents, cam)
	particle_draw.DrawParticles(&entities.Particles, cam)
}
