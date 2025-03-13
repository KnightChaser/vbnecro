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

// VMUser represents a user credential with a role.
type VMUser struct {
	Role     string `yaml:"role"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// VMConfig now includes a list of users.
type VMConfig struct {
	Alias  string   `yaml:"alias"`
	VMName string   `yaml:"vm_name"`
	Users  []VMUser `yaml:"users"`
}

// JobConfig represents a job to perform on a VM.
type JobConfig struct {
	VMAlias    string      `yaml:"vm_alias"`
	EnsureOff  bool        `yaml:"ensure_off,omitempty"`
	Operations []Operation `yaml:"operations"`
}

// Operation holds the type of operation, a role to execute it (if applicable),
// and its parameters.
type Operation struct {
	Type    string                 `yaml:"type"`
	Role    string                 `yaml:"role,omitempty"`
	Params  map[string]interface{} `yaml:"params"`
	StoreAs string                 `yaml:"store_as,omitempty"`
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

// GetUserByRole returns the VMUser with the matching role.
// If not found, it returns an error.
func GetUserByRole(vm *VMConfig, role string) (*VMUser, error) {
	for _, user := range vm.Users {
		if user.Role == role {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with role '%s' not found for VM '%s'", role, vm.VMName)
}
