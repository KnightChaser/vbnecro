package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of our configuration file.
type Config struct {
	VMs  []VMConfig  `yaml:"vms"`
	Jobs []JobConfig `yaml:"jobs"`
}

// VMConfig holds the VirtualBox VM configuration.
type VMConfig struct {
	Alias    string `yaml:"alias"`
	VMName   string `yaml:"vm_name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// JobConfig represents a job to perform on a VM.
type JobConfig struct {
	VMAlias    string      `yaml:"vm_alias"`
	Operations []Operation `yaml:"operations"`
}

// Operation holds the type of operation and its parameters.
// Using Params as a map[string]interface{} allows us to flexibly pass any parameters.
type Operation struct {
	Type   string                 `yaml:"type"`
	Params map[string]interface{} `yaml:"params"`
}

// LoadConfig reads and parses the YAML configuration.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetVMConfig finds a VM configuration by its alias.
func GetVMConfig(vms []VMConfig, alias string) (*VMConfig, error) {
	for _, vm := range vms {
		if vm.Alias == alias {
			return &vm, nil
		}
	}
	return nil, fmt.Errorf("VM with alias '%s' not found", alias)
}
