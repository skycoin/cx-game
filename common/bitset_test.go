package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBitSet(t *testing.T) {
	//arrange
	secondBit := 1

	sixteenthBit := 15
	myBitSet := NewBitSet(16)

	//act
	myBitSet.SetBit(secondBit, true)
	value := myBitSet.GetBit(secondBit)

	//two ticks
	myBitSet.SetBit(sixteenthBit, true)
	myBitSet.SetBit(sixteenthBit, true)
	value2 := myBitSet.GetBit(sixteenthBit)

	//assert

	expected := byte(1)
	assert.Equal(t, expected, value)

	expected = byte(2)
	assert.Equal(t, expected, value2)

}

// func TestGetBitSet(t *testing.T) {
// 	//arrange
// 	secondBit := 1 << 2

// 	//act

// 	//assert

// }
