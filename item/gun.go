package item

import (
	"log"

	"github.com/skycoin/cx-game/spriteloader"
	"github.com/skycoin/cx-game/particles"
)

const bulletSpeed float32 = 40

func UseGun(info ItemUseInfo) {
	log.Print("trying to shoot gun")
	target := info.WorldCoords()
	velocity := target.Sub(info.PlayerCoords()).Normalize().Mul(bulletSpeed)
	particles.CreateBullet( info.PlayerCoords(), velocity )
}

func RegisterGunItemType() ItemTypeID {
	spriteId := spriteloader.LoadSingleSprite(
		"./assets/item/gun-temp.png", "gun" )
	itemType := NewItemType(spriteId)
	itemType.Use = UseGun
	itemType.Name = "Gun"
	return AddItemType(itemType)
}
