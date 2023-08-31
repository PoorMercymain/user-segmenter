package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"RUN_ADDRESS"`
	DatabaseURI   string `env:"DATABASE_URI"`
}

func GetServerConfig() *Config {
	var envCfg Config

	outCfg := getServerFlags()

	_ = env.Parse(&envCfg)

	flag.Parse()

	foundDSNFlag := false
	foundAddressFlag := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "d" {
			foundDSNFlag = true
		}

		if f.Name == "a" {
			foundAddressFlag = true
		}
	})

	if envCfg.ServerAddress != "" && !foundAddressFlag {
		outCfg.ServerAddress = envCfg.ServerAddress
	}

	if envCfg.DatabaseURI != "" && !foundDSNFlag {
		outCfg.DatabaseURI = envCfg.DatabaseURI
	}

	return outCfg
}

func getServerFlags() (cfg *Config) {
	cfg = &Config{}
	cfg.ServerAddress = *flag.String("a", "http://localhost:8080", "server address")
	cfg.DatabaseURI = *flag.String("d", "host=localhost dbname=user-segmenter user=user-segmenter password=user-segmenter port=5432 sslmode=disable", "postgres DSN")
	return
}
