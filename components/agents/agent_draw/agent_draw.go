package agent_draw

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/types"
)

func DrawAgents(agentslist *agents.AgentList, cam *camera.Camera) {
	agentsToDraw := FrustumCull(agentslist.Agents, cam)

	bins := BinByDrawHandlerID(agentsToDraw)

	for drawHandlerID, agents := range bins {
		GetDrawHandler(drawHandlerID)(agents)
	}

}

func FrustumCull(agentlist []*agents.Agent, cam *camera.Camera) []*agents.Agent {
	//todo
	var agentsToDraw []*agents.Agent
	for _, agent := range agentlist {
		if cam.IsInBoundsRadius(agent.PhysicsState.Pos, agent.PhysicsParameters.Radius) {
			agentsToDraw = append(agentsToDraw, agent)
		}
	}
	return agentsToDraw
}

func BinByDrawHandlerID(agentlist []*agents.Agent) map[types.AgentDrawHandlerID][]*agents.Agent {
	bins := make(map[types.AgentDrawHandlerID][]*agents.Agent)

	for _, agent := range agentlist {
		bins[agent.DrawHandlerID] =
			append(bins[agent.DrawHandlerID], agent)
	}

	return bins
}
