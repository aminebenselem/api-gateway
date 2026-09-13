package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type GatewayConfig struct {
	Port            string
	RoutesConfigURI string
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

	viper.AutomaticEnv()

	port := viper.GetString("PORT")
	routesConfigURI := viper.GetString("ROUTES_CONFIG_URI")

	if port == "" {
		return nil, nil, fmt.Errorf("PORT is required")
	}

	if routesConfigURI == "" {
		return nil, nil, fmt.Errorf("ROUTES_CONFIG_URI is required")
	}

	gateway := &GatewayConfig{
		Port:            port,
		RoutesConfigURI: routesConfigURI,
	}

	file, err := os.ReadFile(routesConfigURI)
	if err != nil {
		return nil, nil, fmt.Errorf("read routes config: %w", err)
	}

	var routes RoutesConfig

	if err := json.Unmarshal(file, &routes); err != nil {
		return nil, nil, fmt.Errorf("parse routes config: %w", err)
	}

	return gateway, &routes, nil
}