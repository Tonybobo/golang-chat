package server

import (
	"sync"

	"github.com/tonybobo/go-chat/internal/service"
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

func ConsumerKafkaMsg(data []byte){
	WebSocketServer.BroadCast <- data
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
		case message := <-s.BroadCast:
			msg := &protocol.Message{}
			proto.Unmarshal(message , msg)
			log.Logger.Info("Login", log.Any("Adding User to Online Channel", msg))

			if msg.To != "" {
				//sending msg
				if msg.ContentType >= constant.TEXT && msg.ContentType <= constant.VIDEO {
					//check if sender exist on client map
					_ , exists := s.Clients[msg.From]
					if exists {
						service.MessageService.SaveMessage(msg)
					}
					if msg.MessageType == constant.MESSAGE_TYPE_USER {
						client , online := s.Clients[msg.To]
						if online {
							msgByte , err := proto.Marshal(msg)

							if err == nil {
								client.Send <- msgByte
							}
						}
					}else if msg.MessageType == constant.MESSAGE_TYPE_GROUP {
						sendGroupMessage(msg ,s)
					}
				}else{
					//content type >= 6
					
					client , online := s.Clients[msg.To]
					if online {
						client.Send <- message
					}
				}
			}else {
				//broadcast to all users
				for name ,client :=  range s.Clients {
					log.Logger.Info("name" , log.String("client : ", name ))

					select {
					case client.Send <- message :
					default:
						close(client.Send)
						delete(s.Clients , client.User)
					}
				}
			}
		}
	}
}

func sendGroupMessage(msg *protocol.Message , s *Server) {
	//get all users in the group
	users , _ := service.GroupSerivce.GetGroupUsers(msg.To)

	//get sender details 
	sender := service.UserService.GetUserDetails(msg.From)

	//loop and send 
	for _ , user := range users {
		if user["uid"] == msg.From {
			continue
		}

		recipient , online := s.Clients[user["uid"].(string)]

		if !online {
			continue
		}

		sendMsg := protocol.Message{
			Avatar: sender.Avatar,
			FromUsername: msg.FromUsername,
			From: msg.From,
			To: msg.To,
			Content: msg.Content,
			ContentType: msg.ContentType,
			Type: msg.Type,
			Url: msg.Url,
		}

		msgByte , err := proto.Marshal(&sendMsg)

		if err == nil {
			recipient.Send <- msgByte
		}


	}
}