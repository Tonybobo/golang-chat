package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string `mapstructure:"APPNAME"`
	DBURi   string `mapstructure:"MONGODB_URL"`
	DB      string `mapstructure:"DATABASE"`

	ChannelType           string        `mapstructure:"CHANNEL_TYPE"`
	KafkaHost             string        `mapstructure:"KAFKA_HOST"`
	KafkaTopic            string        `mapstructure:"KAFKA_TOPIC"`
	AccessTokenPrivateKey string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey  string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge     int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	Bucket                string        `mapstructure:"BUCKET"`
	ProjectID             string        `mapstructure:"PROJECT_ID"`
	URL                   string        `mapstructure:"URL"`
	DefaultUserAvatar     string        `mapstructure:"DEFAULT_USER_AVATAR"`
	DefaultGroupAvatar    string        `mapstructure:"DEFAULT_GROUP_AVATAR"`
}

var c Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	viper.Unmarshal(&c)
}

func GetConfig() Config {
	return c
}
