package cxecs

import (
	"github.com/EngoEngine/ecs"
	"github.com/skycoin/cx-game/spriteloader"
)

var world ecs.World

func Init() {
	world = ecs.World{}
	loadAssets()

	//add systems

	//add entities to systems

}

func loadAssets() {
	spriteloader.LoadSingleSprite("./assets/enemies/basic-enemy.png", "enemy")
}

func Update(dt float32) {
	world.Update(dt)
}
