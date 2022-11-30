package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName        string
	MySql          MySqlConfig
	Log            LogConfig
	MsgChannelType MsgChannelType
	Token          Token
}

type MySqlConfig struct {
	HostnPort   string
	Name        string
	Password    string
	TablePrefix string
	User        string
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
