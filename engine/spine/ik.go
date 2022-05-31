package spine

type IKConstraint struct {
	Name string
}

func (constraint *IKConstraint) GetName() string { return constraint.Name }

func (constraint *IKConstraint) Update(skel *Skeleton) {
	// TODO
}
