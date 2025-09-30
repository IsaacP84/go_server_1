package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

type HttpConfig struct {
	Port                  string `mapstructure:"port"`
	Static_file_directory string `mapstructure:"static_file_directory"`
}

type Config struct {
	Server HttpConfig `mapstructure:"server"`
}

var LoadedConfig Config

func LoadConfig(path string) (config Config) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(path)

	// Enable automatic reading of environment variables
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if you want to rely solely on env vars
			fmt.Println("No .env file found, relying on environment variables.")
		} else {
			log.Fatalf("Error reading config file: %v", err)
			panic(err)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
		panic(err)
	}

	p, err := filepath.Abs(config.Server.Static_file_directory)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
	} else {
		config.Server.Static_file_directory = p
	}

	LoadedConfig = config
	fmt.Printf("Config struct: %+v\n", config)

	return config
}
