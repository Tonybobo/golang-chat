package server

import (
	"sync"

	"github.com/tonybobo/go-chat/pkg/common/constant"
	"github.com/tonybobo/go-chat/pkg/global/log"
	"github.com/tonybobo/go-chat/pkg/protocol"
	"google.golang.org/protobuf/proto"
)

var WebSocketServer = NewServer()

type Server struct {
	Clients   map[string]*Client
	mutex     *sync.Mutex
	BroadCast chan []byte
	Online    chan *Client
	Offline   chan *Client
}

func NewServer() *Server {
	return &Server{
		mutex:     &sync.Mutex{},
		Clients:   make(map[string]*Client),
		BroadCast: make(chan []byte),
		Online:    make(chan *Client),
		Offline:   make(chan *Client),
	}
}

func (s *Server) Start() {
	log.Logger.Info("Start Websocket", log.String("starting websocket", "websocket..."))

	for {
		select {
		case conn := <-s.Online:
			log.Logger.Info("Login", log.String("Adding User to Online Channel", conn.User))

			s.Clients[conn.User] = conn
			msg := &protocol.Message{
				Content: constant.PONG,
				Type : constant.HEART_BEAT,
			}

			protoMsg, err := proto.Marshal(msg)
			if err != nil {
				log.Logger.Error("Protocol Error", log.Any("Error:", err))
			}
			conn.Send <- protoMsg
		case conn := <-s.Offline:
			log.Logger.Info("User logging out", log.String("user:", conn.User))
			if _, ok := s.Clients[conn.User]; ok {
				close(conn.Send)
				delete(s.Clients, conn.User)
			}
		}
	}
}
