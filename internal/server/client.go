package server

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/tonybobo/go-chat/pkg/global/log"
	"github.com/tonybobo/go-chat/pkg/protocol"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	Conn *websocket.Conn
	User string
	Send chan []byte
}

func (c *Client) Receiver() {
	defer func() {
		WebSocketServer.Offline <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Logger.Error("Client Read Message Error", log.Any("error:", err))
			WebSocketServer.Offline <- c
			c.Conn.Close()
			break
		}

		msg := &protocol.Message{}
		proto.Unmarshal(message, msg)
		log.Logger.Info("Checking message in development", log.Any("message", msg))
		fmt.Print(msg.Content)
	}
}

func (c *Client) Writer() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		c.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}
