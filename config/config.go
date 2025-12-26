package config

import (
	"io"
	"time"

	"gopkg.in/yaml.v2"
)

// Config contains application configuration
type Config struct {
	ACMEFile        string        `yaml:"acmefile"`
	CertDir         string        `yaml:"certdir"`
	KeyDir          string        `yaml:"keydir"`
	Domains         []string      `yaml:"domains"`
	ExtractInterval time.Duration `yaml:"extractinterval"`
}

// Load decodes a YAML file to give a config structure
func Load(in io.Reader) (*Config, error) {
	var ret Config
	decoder := yaml.NewDecoder(in)
	err := decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}

	// Check if defaults should be used
	if ret.ACMEFile == "" {
		ret.ACMEFile = "/acme.json"
	}
	if ret.CertDir == "" {
		ret.CertDir = "/certs"
	}
	if ret.KeyDir == "" {
		ret.KeyDir = "/private"
	}
	if ret.ExtractInterval < 0 {
		ret.ExtractInterval = 0
	}

	return &ret, nil
}
