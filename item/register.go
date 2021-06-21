package item

var (
	LaserGunItemTypeID  ItemTypeID
	GunItemTypeID ItemTypeID
)

func RegisterItemTypes() {
	LaserGunItemTypeID = RegisterLaserGunItemType()
	GunItemTypeID = RegisterGunItemType()
}
