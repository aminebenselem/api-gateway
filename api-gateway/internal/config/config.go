package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type GatewayConfig struct {
	Port         string
	RoutesConfigURI  string
}

type Route struct {
	Path    string   `json:"path"`
	Targets []string `json:"targets"`
}

type RoutesConfig struct {
	Routes []Route `json:"routes"`
}

func Load() (*GatewayConfig, *RoutesConfig, error) {
	_ = godotenv.Load()

	gateway := &GatewayConfig{
		Port:         os.Getenv("PORT"),
		RoutesConfigURI : os.Getenv("ROUTES_CONFIG_URI"),
	}

	if gateway.Port == "" {
		return nil, nil, fmt.Errorf("PORT is required")
	}

	if gateway.RoutesConfigURI  == "" {
		return nil, nil, fmt.Errorf("ROUTES_CONFIG_URI is required")
	}

	file, err := os.ReadFile(gateway.RoutesConfigURI )
	if err != nil {
		return nil, nil, err
	}

	var routes RoutesConfig

	if err := json.Unmarshal(file, &routes); err != nil {
		return nil, nil, err
	}

	return gateway, &routes, nil
}