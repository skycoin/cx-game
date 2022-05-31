package spine

type CurveTimeline struct {
	Time  []float32
	Curve []Curve
}

func (t *CurveTimeline) Duration() float32 {
	if len(t.Time) == 0 {
		return 0
	}
	return t.Time[len(t.Time)-1]
}

func (t *CurveTimeline) Find(time float32) (int, int, float32) {
	if len(t.Time) == 0 || time < t.Time[0] {
		return -1, -1, 0
	}

	nextKey := 0
	for nextKey < len(t.Time) && time >= t.Time[nextKey] {
		nextKey++
	}

	currentKey := nextKey - 1
	if currentKey < 0 {
		currentKey = 0
	}
	if nextKey >= len(t.Time) {
		nextKey = len(t.Time) - 1
	}

	current := time - t.Time[currentKey]
	delta := t.Time[nextKey] - t.Time[currentKey]
	mix := float32(0.0)
	if delta != 0 {
		mix = current / delta
	}

	if len(t.Curve) > 0 {
		mix = t.Curve[currentKey].Evaluate(mix)
	}
	return currentKey, nextKey, clamp01(mix)
}

type RotateTimeline struct {
	Bone int
	CurveTimeline
	Angle []float32
}

func (t *RotateTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	p0, p1, pmix := t.Find(time)
	bone := skeleton.Bones[t.Bone]
	if p0 < 0 {
		bone.Local.Rotate = lerpAngle(bone.Local.Rotate, 0, alpha)
		return
	}

	target := lerpAngle(t.Angle[p0], t.Angle[p1], pmix)
	bone.Local.Rotate = lerpAngle(bone.Local.Rotate, target, alpha)
}

type TranslateTimeline struct {
	Bone int
	CurveTimeline
	Translate []Vector
}

func (t *TranslateTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	p0, p1, pmix := t.Find(time)
	bone := skeleton.Bones[t.Bone]
	if p0 < 0 {
		bone.Local.Translate = lerpVector(bone.Local.Translate, V0, alpha)
		return
	}

	target := lerpVector(t.Translate[p0], t.Translate[p1], pmix)
	bone.Local.Translate = lerpVector(bone.Local.Translate, target, alpha)
}

type ScaleTimeline struct {
	Bone int
	CurveTimeline
	Scale []Vector
}

func (t *ScaleTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	p0, p1, pmix := t.Find(time)
	bone := skeleton.Bones[t.Bone]
	if p0 < 0 {
		bone.Local.Scale = lerpVector(bone.Local.Scale, V1, alpha)
		return
	}
	target := lerpVector(t.Scale[p0], t.Scale[p1], pmix)
	bone.Local.Scale = lerpVector(bone.Local.Scale, target, alpha)
}

type ShearTimeline struct {
	Bone int
	CurveTimeline
	Shear []Vector
}

func (t *ShearTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	p0, p1, pmix := t.Find(time)
	bone := skeleton.Bones[t.Bone]

	if p0 < 0 {
		bone.Local.Shear = lerpAngleVector(bone.Local.Shear, V0, alpha)
		return
	}

	target := lerpAngleVector(t.Shear[p0], t.Shear[p1], pmix)
	bone.Local.Shear = lerpAngleVector(bone.Local.Shear, target, alpha)
}

type ColorTimeline struct {
	Slot int
	CurveTimeline
	Color []Color
}

func (t *ColorTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	p0, p1, pmix := t.Find(time)
	slot := skeleton.Slots[t.Slot]

	if p0 < 0 {
		slot.Color = lerpColor(slot.Color, slot.Data.Color, alpha)
		return
	}

	target := lerpColor(t.Color[p0], t.Color[p1], pmix)
	slot.Color = lerpColor(slot.Color, target, alpha)
}

type TwoColorTimeline struct {
	Slot int
	CurveTimeline
	Color [][2]Color
}

func (t *TwoColorTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	p0, p1, pmix := t.Find(time)
	slot := skeleton.Slots[t.Slot]

	if p0 < 0 {
		slot.Color = lerpColor(slot.Color, slot.Data.Color, alpha)
		slot.Dark = lerpColor(slot.Dark, slot.Data.Dark, alpha)
		return
	}

	c0, c1 := t.Color[p0], t.Color[p1]
	targetLight := lerpColor(c0[0], c1[0], pmix)
	slot.Color = lerpColor(slot.Color, targetLight, alpha)
	targetDark := lerpColor(c0[1], c1[1], pmix)
	slot.Dark = lerpColor(slot.Dark, targetDark, alpha)
}

type AttachmentTimeline struct {
	Slot int
	CurveTimeline
	Attachment []string
}

func (t *AttachmentTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	if alpha < 0.5 {
		return
	}
	p0, _, _ := t.Find(time)

	slot := skeleton.Slots[t.Slot]
	if p0 < 0 {
		slot.Attachment = skeleton.Attachment(t.Slot, slot.Data.Attachment)
		slot.Deform = nil
		return
	}

	slot.Attachment = skeleton.Attachment(t.Slot, t.Attachment[p0])
	slot.Deform = nil
}

type DeformTimeline struct {
	CurveTimeline
	Skin   *Skin
	Slot   int
	Mesh   string
	Deform [][]Vector
}

func (t *DeformTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	// TODO: figure out proper matchup between attachment and skin
	if t.Skin != skeleton.Skin {
		return
	}

	slot := skeleton.Slots[t.Slot]
	mesh, ok := slot.Attachment.(*MeshAttachment)
	if !ok || t.Mesh != mesh.Name {
		return
	}

	p0, p1, pmix := t.Find(time)
	if p0 < 0 {
		slot.Deform = make([]Vector, len(t.Deform[0]))
		return
	}

	// TODO: this may happen during attachment change
	// this should be reset else where, otherwise we have leftovers
	// from previous attachment deform
	if len(slot.Deform) == 0 {
		slot.Deform = make([]Vector, len(t.Deform[p0]))
	}
	if len(slot.Deform) != len(t.Deform[p0]) {
		panic("invalid state")
	}

	v0, v1 := t.Deform[p0], t.Deform[p1]
	for i := range slot.Deform {
		target := lerpVector(v0[i], v1[i], pmix)
		slot.Deform[i] = lerpVector(slot.Deform[i], target, alpha)
	}
}

type EventTimeline struct {
	CurveTimeline
	Event []string
}

func (t *EventTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	//TODO:

	/*
		if alpha < 0.5 {
			return
		}
		p0, _, _ := t.Find(time)

		var attachment Attachment
		attachmentName := t.Attachment[p0]
		if attachmentName != "" {
			attachment = skeleton.Attachment(t.Slot, attachmentName)
		}

		slot := skeleton.Slots[t.Slot]
		slot.Attachment = attachment
	*/
}

type OrderTimeline struct {
	CurveTimeline
	Order [][]int
}

func (t *OrderTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	if alpha < 0.5 {
		return
	}
	p0, _, _ := t.Find(time)
	if p0 < 0 {
		copy(skeleton.Order, skeleton.Slots)
		return
	}

	order := t.Order[p0]
	if len(order) == 0 {
		for i, slot := range t.Order[p0] {
			skeleton.Order[i] = skeleton.Slots[slot]
		}
	} else {
		copy(skeleton.Order, skeleton.Slots)
	}
}

type TransformConstraintTimeline struct {
	Constraint int
	CurveTimeline
	Transform []TransformMix
}

func (t *TransformConstraintTimeline) Apply(skeleton *Skeleton, time float32, alpha float32) {
	p0, p1, pmix := t.Find(time)
	constraint := skeleton.TransfromConstraints[t.Constraint]

	if p0 < 0 {
		constraint.Mix = LerpTransformMix(constraint.Mix, constraint.Data.Mix, alpha)
		return
	}

	a, b := t.Transform[p0], t.Transform[p1]
	target := LerpTransformMix(a, b, pmix)
	constraint.Mix = LerpTransformMix(constraint.Mix, target, alpha)
}
