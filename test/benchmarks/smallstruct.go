package benchmarks

import (
	"fmt"
	"math/rand"
)

type SmallStruct struct {
	Age  int
	Name string
}

//benchmarking functions
func PassByPointerSmall(smallStruct *SmallStruct) {
	smallStruct.Age += 1
	smallStruct.Name = smallStruct.Name[:len(smallStruct.Name)-1] + fmt.Sprint(rand.Intn(9))
}

func PassByValueSmall(smallStruct SmallStruct) {
	smallStruct.Age += 1
	smallStruct.Name = smallStruct.Name[:len(smallStruct.Name)-1] + fmt.Sprint(rand.Intn(9))
}

func PassByInterfaceSmall(smallStruct interface{}) {
	//conversion
	converted, ok := smallStruct.(SmallStruct)
	if !ok {
		return
	}

	converted.Age += 1
	converted.Name = converted.Name[:len(converted.Name)-1] + fmt.Sprint(rand.Intn(9))
}
