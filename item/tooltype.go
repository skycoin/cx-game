package item

import (
	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
)

func toolTypeFromItemName(itemName string) (types.ToolType,bool) {
	if itemName == "Dev Tile Tool" {
		return constants.TILE_TOOL, true
	}
	if itemName == "Dev Furniture Tool" {
		return constants.FURNITURE_TOOL, true
	}
	if itemName == "Dev Pipe Place Tool" {
		return constants.PIPE_PLACE_TOOL, true
	}
	return types.ToolType(0), false
}
