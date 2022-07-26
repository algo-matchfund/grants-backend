package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config contains the configuration values for the service
type Config struct {
	// Server contains the configuration to run the server on
	Server struct {
		Host           string   `yaml:"host"`
		Port           int      `yaml:"port"`
		Ssl            bool     `yaml:"ssl"`
		Cert           string   `yaml:"cert"`
		Key            string   `yaml:"key"`
		AllowedOrigins []string `yaml:"allowed_origins"`
	} `yaml:"server"`

	// Authentication contains the configuration for the authentication (currently keycloak) server
	Authentication struct {
		PublicKey string `json:"public_key"`
		Host      string `json:"host"`
		Port      int    `json:"port"`
		Realm     string `json:"realm"`
		Client    string `json:"client"`
		Secret    string `json:"secret"`
	} `json:"authentication"`

	// Database contains the configuration needed to connect to the db
	Database struct {
		Port     int    `yaml:"port"`
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`

	// Payment provider configurations for requesting/tracking project donations
	PaymentProviders struct {
		Algorand struct {
			Indexer struct {
				Address string `yaml:"address"`
				Token   string `yaml:"token"`
			} `yaml:"indexer"`
			Node struct {
				Address string `yaml:"address"`
				Token   string `yaml:"token"`
			} `yaml:"node"`
			WatchQueue        int    `yaml:"watch_queue"`
			BlockConfirmation uint64 `yaml:"block_confirmation,omitempty"`
		} `yaml:"algorand"`
	} `yaml:"payment_providers"`

	// SmartContract contains the configuration needed to interact with smart contracts
	SmartContract struct {
		Port  int `yaml:"port"`
		Admin struct {
			Address    string `yaml:"address"`
			Passphrase string `yaml:"passphrase"`
		} `yaml:"admin"`
		Node struct {
			Address string `yaml:"address"`
			Token   string `yaml:"token"`
		} `yaml:"node"`
		BlockConfirmation uint64 `yaml:"block_confirmation,omitempty"`
	} `yaml:"smart_contract"`
}

// LoadConfig loads the configuration from a yaml file
func (c *Config) LoadConfig(filename string) error {
	// Open config file
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&c); err != nil {
		return err
	}

	return nil
}
