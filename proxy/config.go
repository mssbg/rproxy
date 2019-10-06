package proxy

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig() *Proxy {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal("Could not open config.yml")
	}
	defer f.Close()
	yamlBytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read config file %v", err.Error())
	}
	config := Config{Proxy: Proxy{}}

	err = yaml.Unmarshal(yamlBytes, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err.Error())
	}
	config.Proxy.ServiceMap = make(map[string]*Service)
	for _, s := range config.Proxy.Services {
		config.Proxy.ServiceMap[s.Domain] = &s
	}
	config.Proxy.Strategy = GetEnvWithDefault(STRATEGY_VAR_NAME, STRATEGY_RANDOM)
	return &config.Proxy
}

func GetEnvWithDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	} else {
		return value
	}
}
