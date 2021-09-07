package world

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLightSetValues(t *testing.T) {
	myLight := LightValue(253 | 98<<8)

	oldSkyLight := myLight.GetSkyLight() // 128
	oldEnvLight := myLight.GetEnvLight() // 98

	assert.Equal(t, uint8(98), oldSkyLight)
	assert.Equal(t, uint8(253), oldEnvLight)

	myLight.SetSkyLight(25)

	assert.Equal(t, uint8(25), myLight.GetSkyLight())
}
