package render

import (
	"log"
)

var sprites = make(map[string]Sprite)

func RegisterSprite(sprite Sprite) {
	if sprite.Name=="" {
		log.Fatal("cannot register sprite with empty name")
	}
	log.Printf("registering sprite [%v]",sprite.Name)
	sprites[sprite.Name] = sprite
}

func GetSprite(name string) Sprite {
	sprite,ok := sprites[name]
	if !ok { log.Fatalf("cannot get sprite with name [%v]", name) }
	return sprite
}
