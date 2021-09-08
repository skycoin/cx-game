package engine

import (
	"log"
	"os"
	"io/ioutil"
)

var isLogging = true
func ToggleLogging() {
	isLogging = !isLogging
	if isLogging {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}
