package config

import (
	"os"

	"github.com/edgarSucre/mye"
)

type Config struct {
	DbConn    string
	Host      string
	HttpPort  string
	RedisAddr string
	Consumer  string
}

//nolint:errcheck
func LoadConfig(env map[string]string) (Config, error) {
	err := mye.New(mye.CodeInternal, "system_configuration_failed", "missing environment variables")
	for key := range env {
		val := os.Getenv(key)

		if len(val) == 0 {
			err.WithField(key, "value is missing")
		}

		env[key] = val
	}

	if err.HasFields() {
		return Config{}, err
	}

	return Config{
		DbConn:    env["GOOSE_DBSTRING"],
		Host:      env["HTTP_HOST"],
		HttpPort:  env["HTTP_PORT"],
		RedisAddr: env["REDIS_ADDR"],
		Consumer:  env["CONSUMER"],
	}, nil
}
