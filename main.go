package main

import (
	"log"

	appConfig "github.com/andreaswachs/lazyworkflows/appconfig"
)

func main() {
	config := appConfig.New()

	err := config.Load()
	if err != nil {
		log.Fatalf("Could not load config file")
	}

	for repo := range config.Repos {
		log.Printf("I got a new repo: %v\n", repo)
	}
}
