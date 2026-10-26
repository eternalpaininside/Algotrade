package postgrepool

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type PoolConfig struct {
	Host       string        `envconfig:"HOST" required:"true"`
	Port       string        `envconfig:"PORT" default:"5432"`
	PGUser     string        `envconfig:"USER" required:"true"`
	PGPassword string        `envconfig:"PASSWORD" required:"true"`
	PGDatabase string        `envconfig:"DB" required:"true"`
	Timeout    time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfig() (PoolConfig, error) {
	var config PoolConfig

	if err := envconfig.Process(
		"POSTGRES",
		&config,
	); err != nil {
		return PoolConfig{}, fmt.Errorf(
			"process envconfig: %w",
			err,
		)
	}

	return config, nil
}

func NewConfigMust() PoolConfig {
	config, err := NewConfig()

	if err != nil {
		panic(fmt.Errorf(
			"get POSTGRES connection pool config: %w",
			err,
		))
	}

	return config
}
