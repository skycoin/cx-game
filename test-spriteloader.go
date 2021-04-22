package main

import (
	"github.com/skycoin/cx-game/spritesheet"
	"log"
)


func main() {
	log.Print("running test")
	var id int
	id = spritesheet.
		LoadSpriteSheet("./assets/starfield/stars/Starsheet1.png")
	log.Print(id)
	id = spritesheet.
		LoadSpriteSheet("./assets/starfield/stars/Starsheet2.png")
	log.Print(id)
}
