package common

import (
	"gopkg.in/ini.v1"
)

// Config config
type Config struct {
	Common `ini:"common"`
	Server `ini:"server"`
	Redis  `ini:"redis"`
	Log    `ini:"log"`
}

// Common config
type Common struct {
	DataFolder string `ini:"dataFolder"`
	Repo       string `ini:"repo"`
}

// Server config
type Server struct {
	Port       int    `ini:"port"`
	Mode       string `ini:"mode"`
	StaticPath string `ini:"staticPath"`
}

// Server config
type Log struct {
	LogFilePath string `ini:"logFilePath"`
	LogFileName string `ini:"logFileName"`
	LogLevel    int    `ini:"logLevel"`
	LogOutput   string `ini:"logOutput"`
}

// Redis config
type Redis struct {
	Host         string `ini:"host"`
	DB           int    `ini:"db"`
	CachePrefix  string `ini:"cachePrefix"`
	SortedPrefix string `ini:"sortedPrefix"`
}

// InitConfig read config from file
func InitConfig(configFileMame string) (*Config, error) {
	cfg, err := ini.Load(configFileMame)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = cfg.MapTo(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
