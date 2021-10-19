package main

import "fmt"

type One struct {
}

func (o One) PrintOne() {
	fmt.Println("ONE")
}

type Two One

func main() {
	one := One{}
	two := Two{}


	one.PrintOne()
	two.
}
