package item

var (
	LaserGunItemTypeID  ItemTypeID
	GunItemTypeID ItemTypeID
	RockDustItemTypeID ItemTypeID
	BuildToolItemTypeID ItemTypeID
	EnemyToolItemTypeID ItemTypeID
)

func RegisterItemTypes() {
	LaserGunItemTypeID = RegisterLaserGunItemType()
	GunItemTypeID = RegisterGunItemType()
	RockDustItemTypeID = RegisterRockDustItemType()
	BuildToolItemTypeID = RegisterBuildToolItemType()
	EnemyToolItemTypeID = RegisterEnemyToolItemType()

	AddDrops()
}
