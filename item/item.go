package item;

type ItemType struct {
	SpriteID int
}
type ItemTypeID uint32

var itemTypes = []ItemType{}
func NewItemType(SpriteID int) ItemTypeID {
	itemTypes = append(itemTypes,ItemType{SpriteID:SpriteID})
	return ItemTypeID(len(itemTypes)-1)
}
