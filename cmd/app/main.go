package main

import (
	"log"

	"github.com/Akorm0181/yelp/config"
	"github.com/Akorm0181/yelp/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
