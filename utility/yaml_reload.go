package utility

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

//pass in config filename and struct instance
func CheckAndReload(configFilename string, configStruct interface{}) {
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
		if newFileStat.ModTime() != fileStat.ModTime() || newFileStat.Size() != fileStat.Size() {
			data, err := ioutil.ReadFile(configFilename)
			if err != nil {
				log.Panicf("Could not read file: \n, %v", err)
			}
			yaml.Unmarshal(data, configStruct)
			fileStat = newFileStat
		}
		time.Sleep(100 * time.Millisecond)
	}
}
