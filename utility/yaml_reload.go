package utility

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

//pass in config filename and struct instance, and optional channel, run in goroutine
func CheckAndReload(configFilename string, configStruct interface{}, fileHasChanged chan struct{}) {
	var firstCheck bool = true
	// configFilename := "./cmd/starfield/perlin.yaml"
	fileStat, err := os.Stat(configFilename)
	if err != nil {
		log.Panic(err)
	}

	for {
		newFileStat, err := os.Stat(configFilename)
		if err != nil {
			log.Panicf("Could not open file: \n, %v", err)
		}
		//check if file is changed
		if newFileStat.ModTime() != fileStat.ModTime() || newFileStat.Size() != fileStat.Size() || firstCheck {

			data, err := ioutil.ReadFile(configFilename)
			if err != nil {
				log.Panicf("Could not read file: \n, %v", err)
			}
			yaml.Unmarshal(data, configStruct)
			fileStat = newFileStat
			if fileHasChanged != nil {
				fileHasChanged <- struct{}{}
			}
			if !firstCheck {
				fmt.Printf("[File has been changed]: %q\n", configFilename)
			}
			firstCheck = false
		}
		time.Sleep(100 * time.Millisecond)
	}
}
