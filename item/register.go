package item

var (
	LaserGunItemTypeID  ItemTypeID
	GunItemTypeID ItemTypeID
	RockDustItemTypeID ItemTypeID
)

func RegisterItemTypes() {
	LaserGunItemTypeID = RegisterLaserGunItemType()
	GunItemTypeID = RegisterGunItemType()
	RockDustItemTypeID = RegisterRockDustItemType()

	AddDrops()
}
