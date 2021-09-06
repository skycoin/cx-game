package components

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/agents/agent_ai"
	"github.com/skycoin/cx-game/components/agents/agent_health"
	"github.com/skycoin/cx-game/components/agents/agent_physics"
	"github.com/skycoin/cx-game/components/particles/particle_emitter"
	"github.com/skycoin/cx-game/components/particles/particle_physics"
	"github.com/skycoin/cx-game/constants"
)

func Update(dt float32) {
	updateTimers(currentWorld.Entities.Agents.Get(), dt)

	//update lifetimes
	currentWorld.Entities.Particles.Update(dt)

	particle_emitter.Update(dt, currentWorld)
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
	agent_health.UpdateAgents(currentWorld)
	agent_physics.UpdateAgents(currentWorld)
	agent_ai.UpdateAgents(currentWorld, currentPlayer)
	particle_physics.Update(currentWorld)
}
