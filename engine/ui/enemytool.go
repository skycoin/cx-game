package ui

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"
)

type EnemyTool struct {
	scroll int
}

var enemyTool EnemyTool

func mobTypeIDs() []constants.AgentTypeID {
	ids := make([]constants.AgentTypeID,0,constants.NUM_AGENT_TYPES)
	for id := constants.AgentTypeID(0) ; id < constants.NUM_AGENT_TYPES; id++ {
		agenttype := agents.GetAgentType(id)
		if agenttype.Category == constants.AGENT_CATEGORY_FRIENDLY_MOB ||
			agenttype.Category == constants.AGENT_CATEGORY_ENEMY_MOB {
			ids = append(ids, id)
		}
	}
	return ids
}

func DrawEnemyTool(ctx render.Context) {
	enemyTool.Draw(ctx)
}

func EnemyToolScrollUp() { enemyTool.scroll-- }
func EnemyToolScrollDown() { enemyTool.scroll++ }

func (et EnemyTool) Draw(ctx render.Context) {
	et.DrawLine(ctx, et.ActiveAgentID())
}

func (et EnemyTool) DrawLine(
		ctx render.Context, agentTypeID constants.AgentTypeID,
) {
	agentType := agents.GetAgentType(agentTypeID)
	DrawString(agentType.Name, mgl32.Vec4{0,1,0,1}, AlignCenter,ctx)
}

func (et EnemyTool) ActiveAgentID() constants.AgentTypeID {
	mobs := mobTypeIDs()
	idx := cxmath.PositiveModulo(et.scroll, len(mobs))
	return mobs[idx]
}

func EnemyToolActiveAgentID() constants.AgentTypeID {
	return enemyTool.ActiveAgentID()
}
