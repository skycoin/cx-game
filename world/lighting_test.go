package world

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLightSetValues(t *testing.T) {
	myLight := LightValue(15 | 7<<4)

	oldSkyLight := myLight.GetSkyLight() // 7
	oldEnvLight := myLight.GetEnvLight() // 13

	assert.Equal(t, uint8(7), oldSkyLight)
	assert.Equal(t, uint8(15), oldEnvLight)

	myLight.SetSkyLight(8)

	assert.Equal(t, uint8(8), myLight.GetSkyLight())
}
