package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName        string
	Mongo          MongoConfig
	Log            LogConfig
	MsgChannelType MsgChannelType
	Token          Token
	GCP            GcpConfig
}

type MongoConfig struct {
	Username string
	Password string
	Database string
}

type LogConfig struct {
	Path string
}

type Token struct {
	AccessTokenPrivateKey string
	AccessTokenPublicKey  string
	AccessTokenExpiresIn  time.Duration
	AccessTokenMaxAge     int
}

type MsgChannelType struct {
	ChannelType string
	KafkaHost   string
	KafkaTopic  string
}

type GcpConfig struct {
	Bucket        string
	ProjectID     string
	URL           string
	DefaultUserAvatar string
	DefaultGroupAvatar string
}

var c Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

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
