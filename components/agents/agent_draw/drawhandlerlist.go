package agent_draw

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/spriteloader"
)

var DrawHandlerList []func([]*agents.Agent)

func Init() {
	//todo
	spriteloader.LoadSingleSprite("./assets/enemies/basic-enemy.png", "basic-agent")
	RegisterDrawHandler(DrawHandler_1)
	RegisterDrawHandler(DrawHandler_2)
}

func RegisterDrawHandler(newDrawHandler func([]*agents.Agent)) {
	DrawHandlerList = append(DrawHandlerList, newDrawHandler)
}

func ChangeDrawHandler(id int, newDefinition func([]*agents.Agent)) {
	if id < 0 || len(DrawHandlerList) >= id {
		//logging maybe?
		return
	}
	DrawHandlerList[id] = newDefinition
}

func GetDrawHandler(id int) func([]*agents.Agent) {
	return DrawHandlerList[id]
}
