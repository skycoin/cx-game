package agents

import "github.com/skycoin/cx-game/render"

//todo structure better
type PlayerData struct {
	SuitSpriteID         render.SpriteID
	HelmetSpriteID       render.SpriteID
	IgnoringPlatformsFor float32
}
