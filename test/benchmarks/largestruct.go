package benchmarks

import (
	"fmt"
	"math/rand"
)

type LargeStruct struct {
	Age           int
	Name          string
	BirthYear     int
	PassportNuber string
	DriverLicense string
	WorkPlace     string
	innerStruct   InnerStruct
}

type InnerStruct struct {
	InnerField1 int
	InnerField2 string
	InnerFIeld3 float32
}

func PassByInterfaceLarge(largeStruct interface{}) {
	//conversion
	converted, ok := largeStruct.(LargeStruct)
	if !ok {
		return
	}

	converted.Age += 1
	converted.Name = converted.Name[:len(converted.Name)-1] + fmt.Sprint(rand.Intn(9))
}

func PassByPointerLarge(largeStruct *LargeStruct) {
	largeStruct.Age += 1
	largeStruct.Name = largeStruct.Name[:len(largeStruct.Name)-1] + fmt.Sprint(rand.Intn(9))
}

func PassByValueLarge(largeStruct LargeStruct) {
	largeStruct.Age += 1
	largeStruct.Name = largeStruct.Name[:len(largeStruct.Name)-1] + fmt.Sprint(rand.Intn(9))
}
