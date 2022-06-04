package openal_test

import (
	"testing"

	openal "./"
)

func TestGetVendor(t *testing.T) {
	device := openal.OpenDevice("")
	defer device.CloseDevice()

	context := device.CreateContext()
	defer context.Destroy()
	context.Activate()

	vendor := openal.GetVendor()

	if err := openal.Err(); err != nil {
		t.Fatal(err)
	} else if vendor == "" {
		t.Fatal("empty vendor returned")
	}
}
