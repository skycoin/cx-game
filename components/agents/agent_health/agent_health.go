package agent_health

import (
	"github.com/skycoin/cx-game/world"
)

func UpdateAgents(World *world.World) {
	for i, agent := range World.Entities.Agents.Get() {
		if agent.Died() {
			ev :=
				world.NewMobKilledEvent(agent.Meta.Type, World.TimeState.TickCount)
			World.Stats.Log(ev)
			World.Entities.Agents.DestroyAgent(i)
		}
	}
}

func Init() {

}
