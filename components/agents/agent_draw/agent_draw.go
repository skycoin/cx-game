package agent_draw

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/render"
)

func DrawAgents(allAgents *agents.AgentList, cam *camera.Camera) {
	agentsToDraw := FrustumCull(allAgents.Agents, cam)

	bins := BinByDrawHandlerID(agentsToDraw)

	ctx := DrawHandlerContext{Camera: cam}

	for drawHandlerID, agentsForHandler := range bins {
		GetDrawHandler(drawHandlerID)(agentsForHandler, ctx)
	}
	if render.IsBBoxActive() {
		for _, agent := range agentsToDraw {
			render.DrawBBoxLines(agent.Transform.GetInterpolatedBBoxLines(), agent.Transform.GetInterpolatedCollidingLines())
		}
	}

}

func FrustumCull(agentlist []*agents.Agent, cam *camera.Camera) []*agents.Agent {
	//todo
	var agentsToDraw []*agents.Agent
	for _, agent := range agentlist {
		if agent == nil {
			continue
		}
		if cam.IsInBoundsRadius(agent.Transform.Pos, agent.Meta.PhysicsParameters.Radius) {
			agentsToDraw = append(agentsToDraw, agent)
		}
	}
	return agentsToDraw
}

func BinByDrawHandlerID(agentlist []*agents.Agent) map[types.AgentDrawHandlerID][]*agents.Agent {
	bins := make(map[types.AgentDrawHandlerID][]*agents.Agent)

	for _, agent := range agentlist {
		bins[agent.Handlers.Draw] =
			append(bins[agent.Handlers.Draw], agent)
	}

	return bins
}
