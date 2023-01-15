package server

import (
	"github.com/gorilla/websocket"
	"github.com/tonybobo/go-chat/config"
	"github.com/tonybobo/go-chat/pkg/common/constant"
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

		if msg.Type == constant.HEART_BEAT {
			pong := &protocol.Message{
				Content: constant.PONG,
				Type: constant.HEART_BEAT,
			}

			pongByte , err := proto.Marshal(pong)
			if err  != nil {
				log.Logger.Error("Client Marshall Message Error" , log.Any("error : " , err.Error()))
			}

			c.Conn.WriteMessage(websocket.BinaryMessage , pongByte)

		}

		if config.GetConfig().ChannelType == constant.KAKFA {
			//send msg to kafka
			//TODO
		} else {
			WebSocketServer.BroadCast <- message
		}

	}
}

func (c *Client) Writer() {
	defer func() {
		c.Conn.Close()
	}()
	//write message on own websocket
	for message := range c.Send {
		c.Conn.WriteMessage(websocket.BinaryMessage, message)
	}
}
