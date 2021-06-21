package item

import (
	"log"

	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/particles"
)

func UseGun(info ItemUseInfo) {
	log.Print("trying to shoot gun")
	target := info.WorldCoords()
	particles.CreateBullet( info.PlayerCoords(), target )
}

func RegisterGunItemType() ItemTypeID {
	spriteId := spriteloader.LoadSingleSprite(
		"./assets/item/gun-temp.png", "gun" )
	itemType := NewItemType(spriteId)
	itemType.Use = UseGun
	itemType.Name = "Gun"
	return AddItemType(itemType)
}
