package item

import (
	"log"

	"github.com/skycoin/cx-game/spriteloader"
)

func UseGun(info ItemUseInfo) {
	log.Print("trying to shoot gun")
}

func RegisterGunItemType() ItemTypeID {
	spriteId := spriteloader.LoadSingleSprite(
		"./assets/item/gun-temp.png", "gun" )
	itemType := NewItemType(spriteId)
	itemType.Use = UseGun
	itemType.Name = "Gun"
	return AddItemType(itemType)
}
