package components

import (
	"github.com/skycoin/cx-game/camera"
	"github.com/skycoin/cx-game/components/agents/ai"
	"github.com/skycoin/cx-game/components/agents/draw"
	"github.com/skycoin/cx-game/components/agents/health"
	"github.com/skycoin/cx-game/models"
	"github.com/skycoin/cx-game/world"
)

var (
	currentWorldState *world.WorldState
	currentPlanet     *world.Planet
	currentCamera     *camera.Camera
	currentPlayer     *models.Player
)

func Init(planet *world.Planet, cam *camera.Camera, player *models.Player) {
	currentWorldState = planet.WorldState
	currentPlanet = planet
	currentCamera = cam
	currentPlayer = player

	agent_health.Init()
	agent_draw.Init()
	agent_ai.Init()
}

func ChangeCamera(newCamera *camera.Camera) {
	currentCamera = newCamera
}
func ChangePlanet(newPlanet *world.Planet) {
	currentPlanet = newPlanet
	currentWorldState = newPlanet.WorldState
}

func ChangePlayer(newPlayer *models.Player) {
	currentPlayer = newPlayer
}
