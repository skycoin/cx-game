package item

import (
	"github.com/skycoin/cx-game/components/particles/particle_emitter"
	"github.com/skycoin/cx-game/render"
)

const (
	bulletSpeed      float32 = 40
	offsetFromPlayer float32 = 2
)

func UseGun(info ItemUseInfo) {
	target := info.WorldCoords()
	direction := target.Sub(info.PlayerCoords()).Normalize()
	origin := info.PlayerCoords().Add(direction.Mul(offsetFromPlayer))
	velocity := direction.Mul(bulletSpeed)
	particle_emitter.CreateBullet(origin, velocity)
}

func RegisterGunItemType() ItemTypeID {
	/*
		spriteId := spriteloader.LoadSingleSprite(
			"./assets/item/gun-temp.png", "gun")
	*/
	itemType := NewItemType(render.GetSpriteIDByName("gun-temp"))
	itemType.Use = UseGun
	itemType.Name = "Gun"
	return AddItemType(itemType)
}
