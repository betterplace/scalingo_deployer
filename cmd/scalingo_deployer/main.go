package main

import (
	"log"

	sd "github.com/betterplace/scalingo_deployer"
	"github.com/kelseyhightower/envconfig"
)

var config sd.Config

// main is the entry point of the scalingo_deployer application. It initializes
// the configuration from environment variables and starts the deployment
// process. The function uses envconfig to automatically populate the Config
// struct from environment variables. If configuration parsing fails, it
// displays usage information and exits with error code 1.
func main() {
	err := envconfig.Process("", &config)
	if err != nil {
		envconfig.Usage("", &config)
		log.Fatal(err)
	}
	sd.Start(config)
}
