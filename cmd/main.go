package main

import (
	"github.com/tonybobo/go-chat/config"
	"github.com/tonybobo/go-chat/pkg/global/log"
)

func main() {
	log.InitLogger(config.GetConfig().Log.Path)
	log.Logger.Info("config", log.Any("config", config.GetConfig()))

}
