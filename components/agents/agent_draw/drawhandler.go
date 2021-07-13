package agent_draw

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/spriteloader"
)

var DrawHandlerList []func([]*agents.Agent)

func DrawHandler_1(agents []*agents.Agent) {
	for _, agent := range agents {
		spriteloader.DrawSpriteQuad(
			agent.PhysicsState.Pos.X,
			agent.PhysicsState.Pos.Y,
			agent.PhysicsState.Size.X,
			agent.PhysicsState.Size.Y,
			getSpriteId(agent.AgentType),
		)
	}
}

func DrawHandler_2(agents []*agents.Agent) {

}

func Init() {
	spriteloader.LoadSingleSprite("./assets/enemies/basic-enemy.png", "basic-agent")
	DrawHandlerList = append(DrawHandlerList, DrawHandler_1, DrawHandler_2)
}

func AddDrawHandler(newDrawHandler func([]*agents.Agent)) {
	DrawHandlerList = append(DrawHandlerList, newDrawHandler)
}

func GetDrawHandler(id int) func([]*agents.Agent) {
	return DrawHandlerList[id]
}

func getSpriteId(agentType int) spriteloader.SpriteID {
	switch agentType {
	default:
		return spriteloader.GetSpriteIdByName("basic-agent")
	}
}
