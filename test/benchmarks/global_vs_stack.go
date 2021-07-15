package benchmarks

import (
	"fmt"
	"math/rand"
)

type User struct {
	age  int
	name string
}

var user User

func PassByParameter(user User) {
	user.age += 1
	user.name = user.name[:len(user.name)-1] + fmt.Sprint(rand.Intn(5))
}

func PassByGlobal() {
	user.age += 1
	user.name = user.name[:len(user.name)-1] + fmt.Sprint(rand.Intn(5))
}
