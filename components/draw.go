package components

import (
	"github.com/skycoin/cx-game/components/agents/agent_draw"
	"github.com/skycoin/cx-game/components/particles/particle_draw"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/world"
)

//this is to draw the sprites that get "flushed later"
func Draw_Queued(entities *world.Entities, cam *camera.Camera) {
	agent_draw.DrawAgents(&entities.Agents, cam)
}

//draw immediately
func Draw(entities *world.Entities, cam *camera.Camera) {
	particle_draw.DrawParticles(&entities.Particles, cam)
}
