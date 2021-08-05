package components

// TODO replace globals with structs.
// If we have multiple worlds loaded,
// we should be able to simulate them in parallel

import (
	"github.com/skycoin/cx-game/agents"
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/agents/agent_ai"
	"github.com/skycoin/cx-game/components/agents/agent_draw"
	"github.com/skycoin/cx-game/components/agents/agent_health"
	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/components/particles/particle_draw"
	"github.com/skycoin/cx-game/components/particles/particle_physics"
	"github.com/skycoin/cx-game/particle_emitter"
	"github.com/skycoin/cx-game/world"
)

var (
	currentWorld  *world.World
	currentCamera *camera.Camera
	currentPlayer *agents.Agent

	emitter       *particle_emitter.ParticleEmitter
	sparkEmitter  *particle_emitter.SparkEmitter
	bulletEmitter *particle_emitter.BulletEmitter
)

func Init(World *world.World, cam *camera.Camera, player *agents.Agent) {
	/*
		currentWorldState = planet.WorldState
		currentPlanet = planet
		currentCamera = cam
	*/
	currentPlayer = player
	emitter = particle_emitter.
		NewParticle(player.PhysicsState.Pos, &World.Entities.Particles)

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

/*
func ChangePlanet(newPlanet *world.Planet) {
	currentPlanet = newPlanet
	currentWorldState = newPlanet.WorldState
}
*/

func ChangePlayer(newPlayer *agents.Agent) {
	currentPlayer = newPlayer
}
