package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// VMUser represents a user credential with a role.
type VMUser struct {
	Role     string `yaml:"role"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// VMConfig holds the VirtualBox VM configuration.
type VMConfig struct {
	Alias  string   `yaml:"alias"`
	VMName string   `yaml:"vm_name"`
	Users  []VMUser `yaml:"users"`
}

// Operation represents an operation to perform on a VM.
// It includes optional Role and StoreAs fields.
type Operation struct {
	Type    string                 `yaml:"type"`
	Role    string                 `yaml:"role,omitempty"`
	StoreAs string                 `yaml:"store_as,omitempty"`
	Params  map[string]interface{} `yaml:"params"`
}

// JobConfig represents a job to perform on a VM.
type JobConfig struct {
	VMAlias    string      `yaml:"vm_alias"`
	EnsureOff  bool        `yaml:"ensure_off,omitempty"`
	Operations []Operation `yaml:"operations"`
}

// Config represents the complete configuration for the VM manager.
// The vm_manager field is used to flexibly select the backend (e.g. "virtualbox").
type Config struct {
	VMManager string      `yaml:"vm_manager"`
	VMs       []VMConfig  `yaml:"vms"`
	Jobs      []JobConfig `yaml:"jobs"`
}

// loadConfig loads the configuration from the given file path.
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

// GetUserByRole returns the VMUser for the given role from a VMConfig.
func GetUserByRole(vm *VMConfig, role string) (*VMUser, error) {
	for _, user := range vm.Users {
		if user.Role == role {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with role '%s' not found for VM '%s'", role, vm.VMName)
}
