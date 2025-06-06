package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment            string        `mapstructure:"ENVIRONMENT"`
	DBSource               string        `mapstructure:"DB_SOURCE"`
	MigrationURL           string        `mapstructure:"MIGRATION_URL"`
	GRPCServerAddress      string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	WebSocketServerAddress string        `mapstructure:"WEBSOCKET_SERVER_ADDRESS"`
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
