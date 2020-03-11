package configs

import (
	"log"

	"gopkg.in/ini.v1"
)

type Config struct {
	Server       `ini:"server"`
	MySQLConfig  `ini:"mysql"`
	SqliteConfig `ini:"sqlite"`
	OracleConfig `ini:"oracle"`
}

type Server struct {
	Host          string `ini:"host"`
	Port          int    `ini:"port"`
	ViewPattern   string `ini:"viewPattern"`
	StaticPattern string `ini:"staticPattern"`
	Env           string `ini:"env"`
	Dbtype        string `ini:"dbtype"`
}

type MySQLConfig struct {
	Host             string `ini:"host"`
	Port             int    `ini:"port"`
	Dbname           string `ini:"dbname"`
	User             string `ini:"user"`
	Passwd           string `ini:"passwd"`
	Charset          string `ini:"charset"`
	MaxOpenConns     int    `ini:"maxOpenConns"`
	MaxIdleConns     int    `ini:"maxIdleConns"`
	MaxLifetimeConns int    `ini:"maxLifetimeConns"`
}

type SqliteConfig struct {
	Filename string `ini:"filename"`
}

type OracleConfig struct {
	Host             string `ini:"host"`
	Port             int    `ini:"port"`
	Username         string `ini:"username"`
	Password         string `ini:"password"`
	Sid              string `ini:"sid"`
	MaxOpenConns     int    `ini:"maxOpenConns"`
	MaxIdleConns     int    `ini:"maxIdleConns"`
	MaxLifetimeConns int    `ini:"maxLifetimeConns"`
}

// SystemConfig global config
var SystemConfig *Config

func init() {
	var err error
	SystemConfig, err = Load("config.ini")
	if err != nil {
		log.Println(err)
	}
}

// Load get config from ini
func Load(configFileMame string) (*Config, error) {
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
