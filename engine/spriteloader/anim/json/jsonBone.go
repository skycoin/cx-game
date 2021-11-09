package animjson

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

type Bones struct {
}

type Slots struct {
}

type Ik struct {
}

type Skins struct {
}

type Animations struct {
}
