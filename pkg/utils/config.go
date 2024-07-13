package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DNSServerAddr string `mapstructure:"dnsServerAddr"`
	DOHServerAddr string `mapstructure:"dohServerAddr"`
	DoHCertFile   string `mapstructure:"dohCertFile"`
	DoHKeyFile    string `mapstructure:"dohKeyFile"`
}

func LoadConfig(filePath string) Config {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return config
}
