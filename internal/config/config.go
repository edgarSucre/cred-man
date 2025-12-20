package config

import (
	"fmt"

	"github.com/edgarSucre/crm/pkg/terror"
)

type Config struct {
	DbConn    string
	Host      string
	HttpPort  string
	RedisAddr string
	Consumer  string
}

func ErrLoadConfig(key string) error {
	return terror.Internal.New(
		"bad-config",
		fmt.Sprintf("failed to load environment variable %s", key),
	)
}
