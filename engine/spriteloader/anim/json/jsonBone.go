package animjson

import "encoding/json"

type Skeleton struct {
	Hash   string `json:"hash"`
	Spine  string `json:"spine"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Images string `json:"images"`
	Audio  string `json:"audio"`
}

type Bone struct {
	Name     string `json:"name"`
	Parent   string `json:"parent"`
	Length   int    `json:"length"`
	Rotation int    `json:"rotation"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Color    string `json:"color"`
}

type Slot struct {
	Name       string `json:"name"`
	Bone       string `json:"bone"`
	Attachment string `json:"attachment"`
}

type Ikitem struct {
	Name   string            `json:"name"`
	Order  int               `json:"name"`
	Bones  map[string]string `json:"bones"`
	Target string            `json:"target"`
}

type Ik struct {
	Ik map[string]Ikitem `json:"ik"`
}

type Back_arm struct {
	Type   string `json:"type"`
	Hull   int    `json:"hull"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Attachments interface{}

type Skin struct {
	Name        string      `json:"name"`
	Attachments Attachments `json:"attachments"`
}

type AnimatedBoneSpritesheet struct {
	Skeleton   Skeleton    `json:"skeleton"`
	Bones      []Bone      `json:"bones"`
	Slots      []Slot      `json:"slots"`
	Ik         Ik          `json:"ik"`
	Skins      []Skin      `json:"skins"`
	Animations interface{} `json:"animations"`
}

func UnmarshalAnimatedBoneSpriteSheet(buf []byte) AnimatedBoneSpritesheet {
	animatedBoneSpritesheet := AnimatedBoneSpritesheet{}
	json.Unmarshal(buf, &animatedBoneSpritesheet)
	return animatedBoneSpritesheet
}
