package particle_physics

import (
	"fmt"

	"github.com/skycoin/cx-game/components/particles"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/world"
)

func PhysicsHandlerDissappearOnHitCallback(particleList []*particles.Particle, planet *world.Planet) {
	for _, par := range particleList {
		par.MoveNoBounceGravity(planet, constants.TimeStep)
		if par.Collisions.Collided() {
			// fmt.Println("Collidded")
			collisions := ""
			counter := 0
			if par.Collisions.Below {
				counter++
				collisions += "bottom | "
			}
			if par.Collisions.Above {
				counter++
				collisions += "top | "
			}
			if par.Collisions.Left {
				counter++
				collisions += "left | "
			}
			if par.Collisions.Right {
				counter++
				collisions += "right | "
			}
			if counter > 1 {
				fmt.Println("MULTIPLE COLLISIONS!", collisions)
			}
			par.Callback(par)
			par.Die()
		}
	}
}
