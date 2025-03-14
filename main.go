package main

import (
	"flag"
	"log"

	"vnecro/config"
)

func main() {
	// Define command-line flag for configuration file path.
	configPath := flag.String("config-path", "", "Path to the YAML configuration file")
	flag.Parse()
	if *configPath == "" {
		log.Fatal("Missing required flag: --config-path")
	}

	// Load configuration from YAML.
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config from '%s': %v", *configPath, err)
	}

	// Process jobs defined in the config.
	ProcessJobs(cfg)
}
