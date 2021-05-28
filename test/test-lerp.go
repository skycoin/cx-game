package main

import (
	"fmt"
	"time"

	"github.com/skycoin/cx-game/utility"
)

func main() {

	start := float32(5)
	finish := float32(3)
	percent := float32(0)
	for {
		percent += 0.01

		result := utility.Lerp(start, finish, percent)

		fmt.Println(result)
		time.Sleep(time.Millisecond * 50)
	}
}
