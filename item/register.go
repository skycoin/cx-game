package item

var (
	LaserGunItemTypeID  ItemTypeID
	GunItemTypeID ItemTypeID
	RockDustItemTypeID ItemTypeID
	BuildToolItemTypeID ItemTypeID
)

func RegisterItemTypes() {
	LaserGunItemTypeID = RegisterLaserGunItemType()
	GunItemTypeID = RegisterGunItemType()
	RockDustItemTypeID = RegisterRockDustItemType()
	BuildToolItemTypeID = RegisterBuildToolItemType()

	AddDrops()
}
