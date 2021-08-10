package render

import (
	"log"
)

var sprites = []Sprite{}
var spriteNamesToIDs = map[string]int{}

func addSprite(sprite Sprite) int {
	id := len(sprites)
	sprites = append(sprites, sprite)
	return id
}

func RegisterSprite(sprite Sprite) {
	if sprite.Name=="" {
		log.Fatal("cannot register sprite with empty name")
	}
	spriteNamesToIDs[sprite.Name] = addSprite(sprite)
}

func GetSpriteIDByName(name string) SpriteID {
	id,ok := spriteNamesToIDs[name]
	if !ok {
		log.Fatalf("render: cannot find sprite [%v]", name)
	}
	return SpriteID(id)
}

func GetSpriteByID(id SpriteID) Sprite {
	return sprites[id]
}
