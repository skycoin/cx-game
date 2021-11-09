package animjson

import (
	"encoding/json"
)

type Frame struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type SpriteSourceSize struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type SourceSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type KeyFrame struct {
	Frame            Frame            `json:"frame"`
	Rotated          bool             `json:"rotated"`
	Trimmed          bool             `json:"trimmed"`
	SpriteSourceSize SpriteSourceSize `json:"spriteSourceSize"`
	SourceSize       SourceSize       `json:"sourceSize"`
	Duration         int              `json:"duration"`
}

type MetaSize struct {
	W int `json:"w"`
	H int `json:"h"`
}

type FrameTagItem struct {
	Name      string `json:"name"`
	From      int    `json:"from"`
	To        int    `json:"to"`
	Direction string `json:"direction"`
	Color     string `json:"color"`
}

type Meta struct {
	Image     string         `json:"image"`
	Format    string         `json:"format"`
	Size      MetaSize       `json:"size"`
	Scale     string         `json:"scale"`
	FrameTags []FrameTagItem `json:"frameTags"`
}

type AnimatedSpritesheet struct {
	Frames map[string]KeyFrame `json:"frames"`
	Meta   Meta                `json:"meta"`
}

func UnmarshalAnimatedSpriteSheet(buf []byte) AnimatedSpritesheet {
	animatedSpritesheet := AnimatedSpritesheet{}
	json.Unmarshal(buf, &animatedSpritesheet)
	return animatedSpritesheet
}
