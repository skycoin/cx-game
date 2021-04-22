package main

import (
	"github.com/skycoin/cx-game/spriteloader"
	"log"
)


func main() {
	log.Print("running test")
	var id int
	id = spriteloader.
		LoadSpriteSheet("./assets/starfield/stars/Starsheet1.png")
	log.Print(id)
	id = spriteloader.
		LoadSpriteSheet("./assets/starfield/stars/Starsheet2.png")
	log.Print(id)
}
