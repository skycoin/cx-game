package movement

import "testing"

func TestMovementTypeAdded(t *testing.T) {
	movementComponent := NewPlayerMovementComponent()
	movementComponent.AddSupportedMovementType(FLYING)

	if !movementComponent.IsMovementTypePresent(FLYING) {
		t.Error("FLYING SHOULD BE PRESENT!")
	}
}

func TestMovementTypeRemoved(t *testing.T) {
	movementComponent := NewPlayerMovementComponent()
	movementComponent.AddSupportedMovementType(FLYING)

	movementComponent.RemoveSupportedMovementType(FLYING)
	if movementComponent.IsMovementTypePresent(FLYING) {
		t.Error("FLYING SHOULD NOT BE PRESENT")
	}
}

// func TestMovementTypeNotPresent(t *testing.T) {
// 	movementComponent := NewPlayerMovementComponent()

// 	if movementComponent.IsMovementTypePresent(FLYING) {
// 		t.Error("FLYING SHOULD NOT BE PRESENT")
// 	}
// }

func TestMovementTypeChanged(t *testing.T) {
	movementComponent := NewPlayerMovementComponent()

	movementComponent.AddSupportedMovementType(FLYING)
	movementComponent.TryChangeMovementState(FLYING)

	if movementComponent.ActiveMovementType != FLYING {
		t.Error("COULD NOT CHANGE THE STATE")
	}

	movementComponent.RemoveSupportedMovementType(FLYING)

	movementComponent.TryChangeMovementState(FLYING)
	if movementComponent.ActiveMovementType == FLYING {
		t.Error("SHOULD NOT CHANGE THE STATE")
	}
}
