package spine

import (
	"math"
)

type TransformConstraint struct {
	Data *TransformConstraintData

	Bones  []*Bone
	Target *Bone

	Mix TransformMix
}

func NewTransformConstraint(data *TransformConstraintData) *TransformConstraint {
	constraint := &TransformConstraint{}
	constraint.Data = data
	constraint.Mix = data.Mix
	return constraint
}

func (constraint *TransformConstraint) GetName() string {
	return constraint.Data.Name
}

func (constraint *TransformConstraint) Update(skel *Skeleton) {
	if constraint.Data.Local {
		if constraint.Data.Relative {
			constraint.ApplyRelativeLocal()
		} else {
			constraint.ApplyAbsoluteLocal()
		}
	} else {
		if constraint.Data.Relative {
			constraint.ApplyRelativeWorld()
		} else {
			constraint.ApplyAbsoluteWorld()
		}
	}
}

func (constraint *TransformConstraint) ApplyRelativeLocal() {
	panic("todo")
}

func (constraint *TransformConstraint) ApplyAbsoluteLocal() {
	panic("todo")
}

func (constraint *TransformConstraint) ApplyRelativeWorld() {
	panic("todo")
}

func (constraint *TransformConstraint) ApplyAbsoluteWorld() {
	mix := &constraint.Mix
	if mix.Rotate == 0 && mix.Scale <= 0 && mix.Translate == 0 && mix.Shear == 0 {
		return
	}

	target := constraint.Target
	t00, t01, t10, t11 := target.World.M00, target.World.M01, target.World.M10, target.World.M11

	offset := constraint.Data.Offset
	if t00*t11-t01*t10 <= 0 {
		offset.Rotate = -offset.Rotate
		offset.Shear.Y = -offset.Shear.Y
	}

	for _, bone := range constraint.Bones {
		bworld := &bone.World

		if mix.Rotate != 0 {
			b00, b01, b10, b11 := bworld.M00, bworld.M01, bworld.M10, bworld.M11
			r := atan2(t10, t00) - atan2(b10, b00) + offset.Rotate
			if r > math.Pi {
				r -= math.Pi * 2
			} else if r < -math.Pi {
				r += math.Pi * 2
			}
			r *= constraint.Mix.Rotate

			sn, cs := sincos(r)
			bworld.M00 = cs*b00 - sn*b10
			bworld.M01 = cs*b01 - sn*b11
			bworld.M10 = sn*b00 + cs*b10
			bworld.M11 = sn*b01 + cs*b11
		}

		if mix.Translate != 0 {
			t := target.World.Transform(offset.Translate)
			bworld.M02 += (t.X - bworld.M02) * mix.Translate
			bworld.M12 += (t.Y - bworld.M12) * mix.Translate
		}

		if mix.Scale > 0 {
			sx := (sqrt(t00*t00+t10*t10)-1+offset.Scale.X)*constraint.Mix.Scale + 1
			bworld.M00 *= sx
			bworld.M10 *= sx

			sy := (sqrt(t01*t01+t11*t11)-1+offset.Scale.Y)*constraint.Mix.Scale + 1
			bworld.M01 *= sy
			bworld.M11 *= sy
		}

		if mix.Shear > 0 {
			b00, b01, b10, b11 := bworld.M00, bworld.M01, bworld.M10, bworld.M11

			by := atan2(b11, b01)
			r := atan2(t11, t01) - atan2(t10, t00) - (by - atan2(b10, b00))
			if r > math.Pi {
				r -= math.Pi * 2
			} else if r < -math.Pi {
				r += math.Pi * 2
			}
			r = by + (r+offset.Shear.Y)*constraint.Mix.Shear
			s := sqrt(b01*b01 + b11*b11)
			sn, cs := sincos(r)
			bworld.M01 = cs * s
			bworld.M11 = sn * s
		}
	}
}

/*
float rotateMix = this.rotateMix, translateMix = this.translateMix, scaleMix = this.scaleMix, shearMix = this.shearMix;

Bone target = this.target;
float ta = target.a, tb = target.b, tc = target.c, td = target.d;
float degRadReflect = ta * td - tb * tc > 0 ? MathUtils.DegRad : -MathUtils.DegRad;
float offsetRotation = data.offsetRotation * degRadReflect
offsetShearY = data.offsetShearY * degRadReflect;
var bones = this.bones;
for (int i = 0, n = bones.Count; i < n; i++) {
	Bone bone = bones.Items[i];
	bool modified = false;

	if (rotateMix != 0) {
		float a = bone.a, b = bone.b, c = bone.c, d = bone.d;
		float r = MathUtils.atan2(tc, ta) - MathUtils.atan2(c, a) + offsetRotation;
		if (r > MathUtils.PI)
			r -= MathUtils.PI2;
		else if (r < -MathUtils.PI) r += MathUtils.PI2;
		r *= rotateMix;
		float cos = MathUtils.Cos(r), sin = MathUtils.Sin(r);
		bone.a = cos * a - sin * c;
		bone.b = cos * b - sin * d;
		bone.c = sin * a + cos * c;
		bone.d = sin * b + cos * d;
		modified = true;
	}

	if (translateMix != 0) {
		float tx, ty; //Vector2 temp = this.temp;
		target.LocalToWorld(data.offsetX, data.offsetY, out tx, out ty); //target.localToWorld(temp.set(data.offsetX, data.offsetY));
		bone.worldX += (tx - bone.worldX) * translateMix;
		bone.worldY += (ty - bone.worldY) * translateMix;
		modified = true;
	}

	if (scaleMix > 0) {
		float s = (float)Math.sqrt(bone.a * bone.a + bone.c * bone.c);
		//float ts = (float)Math.sqrt(ta * ta + tc * tc);
		if (s > 0.00001f) s = (s + ((float)Math.sqrt(ta * ta + tc * tc) - s + data.offsetScaleX) * scaleMix) / s;
		bone.a *= s;
		bone.c *= s;
		s = (float)Math.sqrt(bone.b * bone.b + bone.d * bone.d);
		//ts = (float)Math.sqrt(tb * tb + td * td);
		if (s > 0.00001f) s = (s + ((float)Math.sqrt(tb * tb + td * td) - s + data.offsetScaleY) * scaleMix) / s;
		bone.b *= s;
		bone.d *= s;
		modified = true;
	}

	if (shearMix > 0) {
		float b = bone.b, d = bone.d;
		float by = MathUtils.atan2(d, b);
		float r = MathUtils.atan2(td, tb) - MathUtils.atan2(tc, ta) - (by - MathUtils.atan2(bone.c, bone.a));
		if (r > MathUtils.PI)
			r -= MathUtils.PI2;
		else if (r < -MathUtils.PI) r += MathUtils.PI2;
		r = by + (r + offsetShearY) * shearMix;
		float s = (float)Math.sqrt(b * b + d * d);
		bone.b = MathUtils.Cos(r) * s;
		bone.d = MathUtils.Sin(r) * s;
		modified = true;
	}

	if (modified) bone.appliedValid = false;
}
*/
