package config

import "gopkg.in/ini.v1"

func LoadConfig(data interface{}, path string) error {
	return ini.MapTo(data, path)
}

// --

type Config struct {
	Port   int16
	Server ServerConfig
	Client ClientConfig
}

type ServerConfig struct {
	Enabled bool
	Host    string
}

type ClientConfig struct {
	Enabled bool
	Host    string
}

func NewConfig() *Config {
	return &Config{
		Port: 4302,
		Server: ServerConfig{
			Enabled: true,
		},
		Client: ClientConfig{
			Enabled: true,
			Host:    "127.0.0.1",
		},
	}
}

func (c *Config) Load(path string) error {
	return LoadConfig(c, path)
}
