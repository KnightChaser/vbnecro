package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"vnecro/config"
)

func main() {
	// Define command-line flag for configuration file path.
	configPath := flag.String("config-path", "", "Path to the YAML configuration file")
	flag.Parse()
	if *configPath == "" {
		logrus.Fatal("Missing required flag: --config-path")
	}

	// Load configuration from YAML.
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logrus.Fatalf("Failed to load config from '%s': %v", *configPath, err)
	}

	// Process jobs defined in the config.
	ProcessJobs(cfg)
}
