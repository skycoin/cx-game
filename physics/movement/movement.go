package movement

type MovementType uint16

type MovementParameters struct {
	MovSpeed         float32
	Mass             float32
	DynamicFriction  float32
	Jumpheight       float32
	Acceleration     float32
	AdditionalJumps  uint
	maxJumps         int
	currentLeftJumps int
}

type MovementComponent struct {
	//active movement state
	ActiveMovementType MovementType
	//active movement state on previous tick
	PreviousActiveMovementType MovementType
	//all movement supported
	movementSet  MovementType
	MovementMeta MovementParameters
	//jump counter

}

func NewPlayerMovementComponent() MovementComponent {
	return MovementComponent{
		ActiveMovementType: NORMAL,
		movementSet:        PLAYER_MOVEMENT_SET,
		MovementMeta: MovementParameters{
			MovSpeed:         DEFAULT_MOVE_SPEED,
			Mass:             DEFAULT_MASS,
			DynamicFriction:  DEFAULT_FRICTION,
			Jumpheight:       DEFAULT_JUMPHEIGHT,
			maxJumps:         2,
			currentLeftJumps: 2,
			Acceleration:     2.0,
		},
	}
}

func (mov *MovementComponent) AddSupportedMovementType(movemenType MovementType) {
	mov.movementSet |= movemenType
}

func (mov *MovementComponent) RemoveSupportedMovementType(movementType MovementType) {
	// if movement type is active, reset it to normal
	if mov.ActiveMovementType == movementType {
		mov.TryChangeMovementState(NORMAL)
	}
	mov.movementSet &= ^movementType
}

func (mov *MovementComponent) IsMovementTypePresent(movementType MovementType) bool {
	return mov.movementSet&movementType != 0
}

func (mov *MovementComponent) TryChangeMovementState(movementType MovementType) bool {
	if mov.IsMovementTypePresent(movementType) {
		mov.PreviousActiveMovementType = mov.ActiveMovementType
		mov.ActiveMovementType = movementType
		return true
	}
	return false
}

func (mov *MovementComponent) ToggleFlying() {
	if mov.ActiveMovementType == FLYING {
		mov.TryChangeMovementState(NORMAL)
	} else {
		mov.TryChangeMovementState(FLYING)
	}
}

func (mov *MovementComponent) ResetJumpCounter() {
	mov.MovementMeta.currentLeftJumps = mov.MovementMeta.maxJumps
}

func (mov *MovementComponent) CanJump() bool {
	//if not standing, or wall sliding
	if mov.ActiveMovementType&(NORMAL|WALL_SLIDING) == 0 {
		return false
	}
	if mov.ActiveMovementType&WALL_SLIDING == 1 && !mov.IsMovementTypePresent(CAN_WALL_JUMP) {
		return false
	}
	if mov.MovementMeta.currentLeftJumps > 0 {
		mov.MovementMeta.currentLeftJumps -= 1
		return true
	}
	return false
}

func (mt MovementType) GetGravityModifier() float32 {
	switch mt {
	case NORMAL:
	case CROUCHING:
		return 1.0
	case FLYING:
		return 0.0
	}
	return 1.0
}

func (mt MovementType) GetFrictionModifier() float32 {
	switch mt {
	case NORMAL:
	case CROUCHING:
		return 1.0
	}
	return 1.0
}

func (mt MovementType) GetMovementSpeedModifier() float32 {
	switch mt {
	case NORMAL:
		return 1
	case CROUCHING:
		return 0.3
	default:
		return 1
	}
}

//for debugging
func (mt MovementType) String() string {
	switch mt {
	case NORMAL:
		return "normal"
	case CROUCHING:
		return "crouching"
	case FLYING:
		return "flying"
	case WALL_SLIDING:
		return "wall_sliding"
	case CLIMBING:
		return "climbing"
	}

	return "unknown"
}
