package agent_health

import (
	"github.com/skycoin/cx-game/world"
)

func UpdateAgents(World *world.World) {
	for _, agent := range World.Entities.Agents.GetAllAgents() {
		if agent == nil {
			continue
		}
		if agent.Died() {
			ev :=
				world.NewMobKilledEvent(agent.Meta.Type, World.TimeState.TickCount)
			World.Stats.Log(ev)
			World.Entities.Agents.DestroyAgent(agent.AgentId)
		}
	}
}

func Init() {

}
