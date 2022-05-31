package spine

import (
	"math"
)

type Bone struct {
	Data   *BoneData
	Parent *Bone
	Local  Transform
	World  Affine
}

func NewBone(data *BoneData) *Bone {
	bone := &Bone{}
	bone.Data = data
	bone.Reset()
	return bone
}

func (bone *Bone) GetName() string { return bone.Data.Name }

func (bone *Bone) Reset() {
	bone.Local = NewTransform()
	bone.World = Identity()
}

func (bone *Bone) Update(skel *Skeleton) {
	var parent *Affine
	if bone.Parent == nil {
		affine := skel.World()
		parent = &affine
	} else {
		parent = &bone.Parent.World
	}
	bone.UpdateWorld(parent, skel.FlipX, skel.FlipY)
}

func (bone *Bone) UpdateWorld(p *Affine, flipX, flipY bool) {
	total := bone.Local.Combine(bone.Data.Local)
	p00, p01, p10, p11 := p.M00, p.M01, p.M10, p.M11
	bone.World.M02 = p00*total.Translate.X + p01*total.Translate.Y + p.M02
	bone.World.M12 = p10*total.Translate.X + p11*total.Translate.Y + p.M12

	switch bone.Data.Inherit {
	case InheritAll:
		snx, csx := sincos(total.Rotate + total.Shear.X)
		sny, csy := sincos(total.Rotate + total.Shear.Y + math.Pi/2)
		t00, t01 := csx*total.Scale.X, csy*total.Scale.Y
		t10, t11 := snx*total.Scale.X, sny*total.Scale.Y

		bone.World.M00 = p00*t00 + p01*t10
		bone.World.M01 = p00*t01 + p01*t11
		bone.World.M10 = p10*t00 + p11*t10
		bone.World.M11 = p10*t01 + p11*t11
		return

	case InheritTranslation:
		snx, csx := sincos(total.Rotate + total.Shear.X)
		sny, csy := sincos(total.Rotate + total.Shear.Y + math.Pi/2)
		t00, t01 := csx*total.Scale.X, csy*total.Scale.Y
		t10, t11 := snx*total.Scale.X, sny*total.Scale.Y

		bone.World.M00 = t00
		bone.World.M01 = t01
		bone.World.M10 = t10
		bone.World.M11 = t11

	case InheritNoRotationOrReflection:
		var prx float32

		s := p00*p00 + p10*p10
		if s > 0.00001 {
			s = abs(p00*p11-p01*p10) / s
			p01 = p10 * s
			p11 = p00 * s
			prx = atan2(p10, p00)
		} else {
			p00 = 0
			p10 = 0
			prx = math.Pi/2 - atan2(p11, p01)
		}

		snx, csx := sincos(total.Rotate + total.Shear.X - prx)
		sny, csy := sincos(total.Rotate + total.Shear.Y - prx + math.Pi/2)
		t00, t01 := csx*total.Scale.X, csy*total.Scale.Y
		t10, t11 := snx*total.Scale.X, sny*total.Scale.Y

		bone.World.M00 = p00*t00 - p01*t10
		bone.World.M01 = p00*t01 - p01*t11
		bone.World.M10 = p10*t00 + p11*t10
		bone.World.M11 = p10*t01 + p11*t11

	case InheritNoScale, InheritNoScaleOrReflection:
		sn, cs := sincos(total.Rotate)
		z00 := p00*cs + p01*sn
		z10 := p10*cs + p11*sn
		s := sqrt(z00*z00 + z10*z10)
		if s > 0.00001 {
			s = 1 / s
		}
		z00 *= s
		z10 *= s

		s = sqrt(z00*z00 + z10*z10)
		r := atan2(z10, z00) + math.Pi/2
		snz, csz := sincos(r)
		z01, z11 := csz*s, snz*s

		snx, csx := sincos(total.Shear.X)
		sny, csy := sincos(total.Shear.Y + math.Pi/2)
		t00, t01 := csx*total.Scale.X, csy*total.Scale.Y
		t10, t11 := snx*total.Scale.X, sny*total.Scale.Y

		flip := false
		if bone.Data.Inherit != InheritNoScaleOrReflection {
			flip = p00*p11-p01*p10 < 0
		} else {
			flip = flipX != flipY
		}
		if flip {
			z01 = -z01
			z11 = -z11
		}

		bone.World.M00 = z00*t00 + z01*t10
		bone.World.M01 = z00*t01 + z01*t11
		bone.World.M10 = z10*t00 + z11*t10
		bone.World.M11 = z10*t01 + z11*t11
		return
	}

	if flipX {
		bone.World.M00 = -bone.World.M00
		bone.World.M01 = -bone.World.M01
	}
	if flipY {
		bone.World.M10 = -bone.World.M10
		bone.World.M11 = -bone.World.M11
	}
}
