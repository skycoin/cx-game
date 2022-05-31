package spine

type PathConstraint struct {
	Name string
}

func (constraint *PathConstraint) GetName() string { return constraint.Name }

func (constraint *PathConstraint) Update(skel *Skeleton) {
	// TODO
}
