package main

import (
	"fmt"
	"log"
	"os"

	appConfig "github.com/andreaswachs/lazyworkflows/appconfig"
)

func main() {
	config := appConfig.New()

	err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not load config file. See error msg.")
		os.Exit(0)
	}

	for repo := range config.Repos {
		log.Printf("I got a new repo: %v\n", repo)
	}
}
