package main

import (
	"log"

	sd "github.com/betterplace/scalingo_deployer"
	"github.com/kelseyhightower/envconfig"
)

var config sd.Config

func main() {
	err := envconfig.Process("", &config)
	if err != nil {
		envconfig.Usage("", &config)
		log.Fatal(err)
	}
	sd.Start(config)
}
