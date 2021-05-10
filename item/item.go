package item;

type ItemType struct {
	SpriteID int
}

var itemTypes = []ItemType{}
func NewItemType(SpriteID int) uint32 {
	itemTypes = append(itemTypes,ItemType{SpriteID:SpriteID})
	return uint32(len(itemTypes)-1)
}
