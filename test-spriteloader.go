package main

import (
	"github.com/skycoin/cx-game/spriteloader"
	//"github.com/skycoin/cx-game/render"
	"log"
)


func main() {
	log.Print("running test")
	spriteSheetId := spriteloader.
		LoadSpriteSheet("./assets/starfield/stars/Starsheet1.png")
	spriteloader.
		LoadSprite(spriteSheetId, "star", 0,0)
	spriteId := spriteloader.
		GetSpriteIdByName("star")
	log.Print(spriteId)
}
