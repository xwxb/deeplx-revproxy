package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
	Default   Default    `yaml:"default"`
	General   General    `yaml:"general"`
	Server    Server     `yaml:"server"`
}

type Endpoint struct {
	URL       string `yaml:"url"`
	Weight    int    `yaml:"weight"`
	Timeout   int    `yaml:"timeout"`
	AuthToken string `yaml:"auth_token"`
}

type Default struct {
	Timeout int `yaml:"timeout"`
}

type General struct {
	MaxRetries int `yaml:"max_retries"`
}

type Server struct {
	Port int `yaml:"port"`
}

var Global Config

func InitConfig() {
	viper.SetConfigFile("./config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := viper.Unmarshal(&Global); err != nil {
		log.Fatal("Failed to unmarshal config: %v", err)
	}

	// Set the default values
	defaultTimeout := viper.GetInt("default.timeout")
	for i := range Global.Endpoints {
		if Global.Endpoints[i].Timeout == 0 {
			Global.Endpoints[i].Timeout = defaultTimeout
		}
	}
}

func LogConfigedEndpoints() {
	log.Println("Configured endpoints:")
	for _, ep := range Global.Endpoints {
		log.Printf("URL: %s, Weight: %d, Timeout: %d\n", ep.URL, ep.Weight, ep.Timeout)
	}
}
