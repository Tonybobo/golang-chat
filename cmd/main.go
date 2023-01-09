package main

import (
	"net/http"
	"time"

	"github.com/tonybobo/go-chat/internal/router"
	"github.com/tonybobo/go-chat/internal/server"
	"github.com/tonybobo/go-chat/pkg/global/log"
)

func main() {
	log.InitLogger("logs/chat.log")

	log.Logger.Info("start server", log.String("start", "server starting"))

	router := router.NewRouter()

	go server.WebSocketServer.Start()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Logger.Error("Fail to Start Server", log.Any("Error: ", err))
	}
}
