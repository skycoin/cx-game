package components

// TODO replace globals with structs.
// If we have multiple worlds loaded,
// we should be able to simulate them in parallel

import (
	"github.com/skycoin/cx-game/components/agents"
	"github.com/skycoin/cx-game/components/agents/agent_ai"
	"github.com/skycoin/cx-game/components/agents/agent_draw"
	"github.com/skycoin/cx-game/components/agents/agent_health"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/particles/particle_draw"
	"github.com/skycoin/cx-game/components/particles/particle_emitter"
	"github.com/skycoin/cx-game/components/particles/particle_physics"
	"github.com/skycoin/cx-game/engine/camera"
	"github.com/skycoin/cx-game/world"
)

var (
	currentWorld  *world.World
	currentCamera *camera.Camera
	currentPlayer *agents.Agent
)

func Init(World *world.World, cam *camera.Camera, player *agents.Agent) {
	currentPlayer = player

	agent_health.Init()
	agent_draw.Init()
	agent_ai.Init()

	particle_physics.Init()
	particle_draw.Init()

	particles.Init()
	particle_emitter.Init(&World.Entities.Particles)

}

func ChangeCamera(newCamera *camera.Camera) {
	currentCamera = newCamera
}

func ChangeWorld(newWorld *world.World) {
	currentWorld = newWorld
}

func ChangePlayer(newPlayer *agents.Agent) {
	currentPlayer = newPlayer
}
