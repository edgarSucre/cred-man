package config

import (
	"fmt"
	"os"

	"github.com/edgarSucre/crm/pkg/terror"
	"github.com/joho/godotenv"
)

type Config struct {
	DbConn   string
	Host     string
	HttpPort string
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	env := map[string]string{
		"GOOSE_DBSTRING": "",
		"HTTP_HOST":      "",
		"HTTP_PORT":      "",
	}

	for key := range env {
		val := os.Getenv(key)

		if len(val) == 0 {
			return Config{}, errLoadConfig(key)
		}

		env[key] = val
	}

	return Config{
		DbConn:   env["GOOSE_DBSTRING"],
		Host:     env["HTTP_HOST"],
		HttpPort: env["HTTP_PORT"],
	}, nil
}

func errLoadConfig(key string) error {
	return terror.Internal.New(
		"bad-config",
		fmt.Sprintf("failed to load environment variable %s", key),
	)
}
