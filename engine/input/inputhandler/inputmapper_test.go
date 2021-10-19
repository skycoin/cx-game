package inputhandler

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/stretchr/testify/assert"
)

func TestInputMapper(t *testing.T) {
	//arrange
	myInputMapper := NewAgentInputMapper()

	//act
	key := myInputMapper.GetKey(MOVE_LEFT)

	//assert
	expectedMoveLeft := glfw.KeyA
	assert.Equal(t, expectedMoveLeft, key)

}

func TestRebindKey(t *testing.T) {

	//arrange
	myInputMapper := NewAgentInputMapper()

	//act
	leftKey := myInputMapper.GetKey(MOVE_LEFT)
	myInputMapper.BindKey(glfw.KeyR, MOVE_LEFT)

	newLeftKey := myInputMapper.GetKey(MOVE_LEFT)

	//assert
	expected1 := glfw.KeyA
	expected2 := glfw.KeyR

	assert.Equal(t, expected1, leftKey)
	assert.Equal(t, expected2, newLeftKey)
}

func TestRebindResetKey(t *testing.T) {
	//arrange
	myInputMapper := NewAgentInputMapper()

	//act
	// leftKey := myInputMapper.GetKey(MOVE_LEFT)
	//rebind a key to new enum
	myInputMapper.BindKey(glfw.KeyA, 15)

	newLeftKey := myInputMapper.GetKey(MOVE_LEFT)

	//assert
	expected := glfw.KeyUnknown

	//expect left move key to be unknown because we rebind it to other enum
	assert.Equal(t, expected, newLeftKey)
	// assert.Equal(t, expected2, newLeftKey)

}
