package spine

type Updateable interface {
	Update(skeleton *Skeleton)
}

type Skeleton struct {
	Data  *SkeletonData
	Local Transform
	FlipX bool
	FlipY bool

	Bones []*Bone
	Slots []*Slot
	Order []*Slot
	Skin  *Skin

	ResetBones  []*Bone
	UpdateOrder []Updateable

	IKConstraints        []*IKConstraint
	TransfromConstraints []*TransformConstraint
	PathConstraints      []*PathConstraint
}

func NewSkeleton(data *SkeletonData) *Skeleton {
	skel := &Skeleton{}
	skel.Data = data
	skel.Skin = data.DefaultSkin
	skel.Local = NewTransform()

	for _, bonedata := range data.Bones {
		skel.Bones = append(skel.Bones, NewBone(bonedata))
	}

	for _, bone := range skel.Bones {
		if bone.Data.Parent != nil {
			bone.Parent = skel.Bones[bone.Data.Parent.Index]
		}
	}

	for _, slotdata := range data.Slots {
		slot := NewSlot(skel.Bones[slotdata.Bone.Index], slotdata)
		skel.Slots = append(skel.Slots, slot)
		skel.Order = append(skel.Order, slot)
	}

	for _, constraintdata := range data.TransformConstraints {
		constraint := NewTransformConstraint(constraintdata)
		skel.TransfromConstraints = append(skel.TransfromConstraints, constraint)
		constraint.Target = skel.Bones[constraintdata.Target.Index]
		constraint.Bones = make([]*Bone, len(constraintdata.Bones))
		for i, bonedata := range constraintdata.Bones {
			constraint.Bones[i] = skel.Bones[bonedata.Index]
		}
	}

	skel.ResetUpdateOrder()

	skel.UpdateAttachments()
	skel.Update()

	return skel
}

func (skel *Skeleton) UpdateAttachments() {
	for _, slot := range skel.Slots {
		slot.Attachment = skel.Attachment(slot.Data.Index, slot.Data.Attachment)
	}
}

func (skel *Skeleton) FindBone(name string) *Bone {
	data := skel.Data.FindBone(name)
	if data == nil {
		return nil
	}
	return skel.Bones[data.Index]
}

func (skel *Skeleton) FindSlot(name string) *Slot {
	data := skel.Data.FindSlot(name)
	if data == nil {
		return nil
	}
	return skel.Slots[data.Index]
}

func (skel *Skeleton) FindTransformConstraint(name string) *TransformConstraint {
	data := skel.Data.FindTransformConstraint(name)
	if data == nil {
		return nil
	}
	return skel.TransfromConstraints[data.Index]
}

func (skel *Skeleton) FindSkin(name string) *Skin {
	return skel.Data.FindSkin(name)
}

func (skel *Skeleton) FindAnimation(name string) *Animation {
	return skel.Data.FindAnimation(name)
}

func (skel *Skeleton) Attachment(slot int, name string) Attachment {
	attachment := skel.Skin.Attachment(slot, name)
	if attachment != nil {
		return attachment
	}
	if skel.Data.DefaultSkin != skel.Skin {
		return skel.Data.DefaultSkin.Attachment(slot, name)
	}
	return nil
}

func (skel *Skeleton) ResetUpdateOrder() {
	skel.ResetBones = []*Bone{}
	skel.UpdateOrder = []Updateable{}

	updateOrderContains := func(v Updateable) bool {
		for _, t := range skel.UpdateOrder {
			if t == v {
				return true
			}
		}
		return false
	}

	boneSorted := make([]bool, len(skel.Bones))
	total := 0
	total += len(skel.IKConstraints)
	total += len(skel.TransfromConstraints)
	total += len(skel.PathConstraints)

	var sortBone func(bone *Bone)
	sortBone = func(bone *Bone) {
		if bone == nil {
			return
		}
		if boneSorted[bone.Data.Index] {
			return
		}
		if bone.Parent != nil {
			sortBone(bone.Parent)
		}
		boneSorted[bone.Data.Index] = true
		skel.UpdateOrder = append(skel.UpdateOrder, bone)
	}

	var resetBones func(bones []*BoneData)
	resetBones = func(bones []*BoneData) {
		for _, bone := range bones {
			if !boneSorted[bone.Index] {
				continue
			}
			resetBones(bone.Children)
			boneSorted[bone.Index] = false
		}
	}

	sortTransform := func(constraint *TransformConstraint) {
		sortBone(constraint.Target)
		if constraint.Data.Local {
			for _, bone := range constraint.Bones {
				sortBone(bone.Parent)
				if !updateOrderContains(bone) {
					skel.ResetBones = append(skel.ResetBones, bone)
				}
			}
		} else {
			for _, bone := range constraint.Bones {
				sortBone(bone)
			}
		}

		skel.UpdateOrder = append(skel.UpdateOrder, constraint)

		for _, bone := range constraint.Bones {
			resetBones(bone.Data.Children)
		}
		for _, bone := range constraint.Bones {
			boneSorted[bone.Data.Index] = true
		}
	}

	for order := 0; order < total; order++ {
		for _, constraint := range skel.TransfromConstraints {
			if constraint.Data.Order == order {
				sortTransform(constraint)
				break
			}
		}
	}

	for _, bone := range skel.Bones {
		sortBone(bone)
	}
}

func (skel *Skeleton) World() Affine {
	root := skel.Local.Affine()
	if skel.FlipX {
		root.M00 = -root.M00
		root.M01 = -root.M01
	}
	if skel.FlipY {
		root.M10 = -root.M10
		root.M11 = -root.M11
	}
	return root
}

func (skel *Skeleton) SetToSetupPose() {
	for _, bone := range skel.Bones {
		bone.Local = NewTransform()
		bone.World = bone.Local.Affine()
	}

	for _, constraint := range skel.TransfromConstraints {
		constraint.Mix = constraint.Data.Mix
	}

	for _, slot := range skel.Slots {
		slot.Attachment = skel.Attachment(slot.Data.Index, slot.Data.Attachment)
		slot.Deform = nil
		slot.Color = slot.Data.Color
		slot.Dark = slot.Data.Dark
	}
}

func (skel *Skeleton) Update() {
	for _, bone := range skel.ResetBones {
		bone.World = bone.Local.Affine()
	}
	for _, update := range skel.UpdateOrder {
		update.Update(skel)
	}
}

type Slot struct {
	Data       *SlotData
	Bone       *Bone
	Attachment Attachment
	Deform     []Vector
	Color      Color
	Dark       Color
}

func NewSlot(bone *Bone, data *SlotData) *Slot {
	slot := &Slot{}
	slot.Bone = bone
	slot.Color = data.Color
	slot.Data = data
	return slot
}
