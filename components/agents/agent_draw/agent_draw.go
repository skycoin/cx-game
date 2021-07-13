package agent_draw

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/camera"
)

func DrawAgents(agentslist *agents.AgentList, cam *camera.Camera) {
	agentsToDraw := FrustumCull(agentslist.Agents, cam)

	sortedByDrawType := SortByType(agentsToDraw)

	for k, v := range sortedByDrawType {
		DrawFuncsArray[k](v)
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

func SortByType(agentlist []*agents.Agent) map[int][]*agents.Agent {
	sortedMap := make(map[int][]*agents.Agent)

	for _, agent := range agentlist {
		sortedMap[agent.DrawFunctionUpdateId] = append(sortedMap[agent.DrawFunctionUpdateId], agent)
	}

	return sortedMap
}
