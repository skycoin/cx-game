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

type Bone struct {
	Name     string `json:"name"`
	Parent   string `json:"parent"`
	Length   int    `json:"length"`
	Rotation int    `json:"rotation"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Color    string `json:"color"`
}

type Bones struct {
	Bones map[string]Bone `json:"bones"`
}

type Slot struct {
	Name       string `json:"name"`
	Bone       string `json:"bone"`
	Attachment string `json:"attachment"`
}

type Slots struct {
	Slots map[string]Slot `json:"slots"`
}

type Ik struct {
}

type Skins struct {
}

type Animations struct {
}
