package kafka

import (
	"strings"

	"github.com/Shopify/sarama"
	"github.com/tonybobo/go-chat/pkg/global/log"
)

var producer sarama.AsyncProducer
var topic string = "default_message"

func InitProducer(topicInput , hosts string){
	topic = topicInput
	config := sarama.NewConfig()
	config.Producer.Compression = sarama.CompressionGZIP
	client , err := sarama.NewClient(strings.Split(hosts , ",") , config)
	if err != nil {
		log.Logger.Error("fail to initialize producer" , log.Any("error" , err.Error()))
	}

	producer , err = sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		log.Logger.Error("fail to initialize producer" , log.Any("error" , err.Error()))
	}

}

func Send(data []byte){
	buffer := sarama.ByteEncoder(data)
	producer.Input() <- &sarama.ProducerMessage{Topic: topic , Key: nil , Value: buffer}
}

func Close(){
	if producer != nil {
		producer.Close()
	}
}