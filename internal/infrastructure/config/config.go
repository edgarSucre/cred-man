package config

import (
	"fmt"
	"os"

	"github.com/edgarSucre/crm/pkg/terror"
)

type Config struct {
	DbConn    string
	Host      string
	HttpPort  string
	RedisAddr string
	Consumer  string
}

func LoadConfig(env map[string]string) (Config, error) {
	for key := range env {
		val := os.Getenv(key)

		if len(val) == 0 {
			return Config{}, ErrLoadConfig(key)
		}

		env[key] = val
	}

	return Config{
		DbConn:    env["GOOSE_DBSTRING"],
		Host:      env["HTTP_HOST"],
		HttpPort:  env["HTTP_PORT"],
		RedisAddr: env["REDIS_ADDR"],
		Consumer:  env["CONSUMER"],
	}, nil
}

func ErrLoadConfig(key string) error {
	return terror.Internal.New(
		"bad-config",
		fmt.Sprintf("failed to load environment variable %s", key),
	)
}
