package components

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/components/agents/agent_ai"
	"github.com/skycoin/cx-game/components/agents/agent_draw"
	"github.com/skycoin/cx-game/components/agents/agent_health"
	"github.com/skycoin/cx-game/components/agents/agent_physics"
	"github.com/skycoin/cx-game/components/particles/particle_physics"
	"github.com/skycoin/cx-game/world"
)

func Update(dt float32) {
	particle_physics.Update(dt)
	updateTimers(currentWorldState.AgentList.Agents, dt)
}

func updateTimers(agents []*agents.Agent, dt float32) {
	for _,agent := range agents {
		if agent.DrawHandlerID == constants.DRAW_HANDLER_ANIM {
			agent.AnimationState.Action.Update(dt)
		}
		if agent.WaitingFor > 0 { agent.WaitingFor -= dt }
		if agent.Died() { agent.TimeSinceDeath += dt }
	}
}

func FixedUpdate() {
	//update health state first
	agent_health.UpdateAgents(currentWorldState.AgentList)
	//update physics state second
	agent_physics.UpdateAgents(currentWorldState, currentPlanet)

	agent_ai.UpdateAgents(currentWorldState.AgentList, currentPlayer)
}

func Draw(worldState *world.WorldState, cam *camera.Camera) {
	agent_draw.DrawAgents(worldState.AgentList, cam)
}
