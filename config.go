package main
import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Profiles map[string]Profile `json:"profiles"`
	Dirs     Dirs               `json:"dirs"`
	Services Services           `json:"services"`
}

type Profile struct {
	AccountPath  string `json:"accountPath"`
	KeyPath      string `json:"keyPath"`
	IdentityPath string `json:"identityPath"`
}

type Dirs struct {
	Identity string `json:"identity"`
	Account  string `json:"account"`
	Key      string `json:"key"`
	Proofs	 string	`json:"proofs"`
}

type RPCProviders struct {
	Local   map[string]string `json:"local"`
	Public  map[string]string `json:"public"`
	Private map[string]string `json:"private"`
}

type Services struct {
	STORAGE           string       `json:"STORAGE"`
	CONTRACT_ADDR     string       `json:"CONTRACT_ADDR"`
	RPC_PROVIDERS_URLS RPCProviders `json:"RPC_PROVIDERS_URLS"`
}

func (c *Config) Load() error {
	data, err := os.ReadFile(".config.json")
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config JSON: %w", err)
	}
	return nil
}

func (c *Config) AddProfile(username, accountPath, keyPath, identityPath string) error {
	// Load existing config (if this func isn't called after Load already)
	if err := c.Load(); err != nil {
		return err
	}

	if c.Profiles == nil {
		c.Profiles = make(map[string]Profile)
	}

	c.Profiles[username] = Profile{
		AccountPath:  accountPath,
		KeyPath:      keyPath,
		IdentityPath: identityPath,
	}

	// Re-save updated JSON
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %w", err)
	}

	if err := os.WriteFile(".config.json", data, 0644); err != nil {
		return fmt.Errorf("failed to write updated config: %w", err)
	}

	return nil
}
