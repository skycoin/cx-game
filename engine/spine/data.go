package spine

type SkeletonData struct {
	Name   string
	Images string
	Size   Vector

	Bones       []*BoneData
	Slots       []*SlotData
	Skins       []*Skin
	DefaultSkin *Skin

	Events     []*EventData
	Animations []*Animation

	// IKConstraints        []*IKConstraintData
	TransformConstraints []*TransformConstraintData
	// PathConstraints      []*PathConstraintData
}

func (skel *SkeletonData) FindBone(name string) *BoneData {
	for _, bone := range skel.Bones {
		if bone.Name == name {
			return bone
		}
	}
	return nil
}

func (skel *SkeletonData) FindSlot(name string) *SlotData {
	for _, slot := range skel.Slots {
		if slot.Name == name {
			return slot
		}
	}
	return nil
}

func (skel *SkeletonData) FindSkin(name string) *Skin {
	for _, skin := range skel.Skins {
		if skin.Name == name {
			return skin
		}
	}
	return nil
}

func (skel *SkeletonData) FindAnimation(name string) *Animation {
	for _, anim := range skel.Animations {
		if anim.Name == name {
			return anim
		}
	}
	return nil
}

func (skel *SkeletonData) FindTransformConstraint(name string) *TransformConstraintData {
	for _, transform := range skel.TransformConstraints {
		if transform.Name == name {
			return transform
		}
	}
	return nil
}

type BoneData struct {
	Index    int
	Name     string
	Length   float32
	Parent   *BoneData
	Local    Transform
	Inherit  Inherit
	Color    Color
	Children []*BoneData
}

type Inherit byte

const (
	InheritTranslation = Inherit(1 << iota)
	InheritRotation
	InheritScale
	InheritReflection

	InheritAll                    = InheritTranslation | InheritRotation | InheritScale | InheritReflection
	InheritNoRotationOrReflection = InheritAll &^ (InheritRotation | InheritReflection)
	InheritNoScale                = InheritAll &^ InheritScale
	InheritNoScaleOrReflection    = InheritAll &^ (InheritScale | InheritReflection)
)

func parseInherit(s string) Inherit {
	switch s {
	default:
		fallthrough
	case "normal", "":
		return InheritAll
	case "onlyTranslation":
		return InheritTranslation
	case "noRotationOrReflection":
		return InheritNoRotationOrReflection
	case "noScale":
		return InheritNoScale
	case "noScaleOrReflection":
		return InheritNoScaleOrReflection
	}
}

type SlotData struct {
	Index      int
	Name       string
	Bone       *BoneData
	Attachment string
	Color      Color
	Dark       Color
	Blend      BlendMode
}

type EventData struct {
	Index  int
	Name   string
	Int    int
	Float  float32
	String string
}

type TransformConstraintData struct {
	Index int

	Name  string
	Order int

	Bones  []*BoneData
	Target *BoneData

	Mix    TransformMix
	Offset Transform

	Relative bool
	Local    bool
}

type TransformMix struct {
	Rotate    float32
	Translate float32
	Scale     float32
	Shear     float32
}

func LerpTransformMix(a, b TransformMix, p float32) TransformMix {
	r := TransformMix{}
	r.Rotate = lerp(a.Rotate, b.Rotate, p)
	r.Translate = lerp(a.Translate, b.Translate, p)
	r.Scale = lerp(a.Scale, b.Scale, p)
	r.Shear = lerp(a.Shear, b.Shear, p)
	return r
}

type Skin struct {
	Name        string
	Attachments map[SkinSlot]Attachment
}

func NewSkin(name string) *Skin {
	return &Skin{
		Name:        name,
		Attachments: map[SkinSlot]Attachment{},
	}
}

type SkinSlot struct {
	SlotIndex  int
	Attachment string
}

func (skin *Skin) Attachment(slot int, attachmentName string) Attachment {
	attach, ok := skin.Attachments[SkinSlot{slot, attachmentName}]
	if !ok {
		return nil
	}
	return attach
}

func (skin *Skin) AddAttachment(slot int, attachmentName string, attach Attachment) {
	skin.Attachments[SkinSlot{slot, attachmentName}] = attach
}

type BlendMode byte

const (
	Normal = BlendMode(iota)
	Additive
	Multiply
	Screen
)

func parseBlendMode(s string) BlendMode {
	switch s {
	case "", "normal":
		return Normal
	case "additive":
		return Additive
	case "multiply":
		return Multiply
	case "screen":
		return Screen
	default:
		return Normal
	}
}
