package kafka

import (
	"strings"

	"github.com/Shopify/sarama"
	"github.com/tonybobo/go-chat/pkg/global/log"
)

var consumer sarama.Consumer

type ConsumerCallback func(data []byte)

func InitConsumer(hosts string){
	config := sarama.NewConfig()
	client , err := sarama.NewClient(strings.Split(hosts , ","),config)
	if err != nil {
		log.Logger.Error("fail to initialize kafka" , log.Any("error" , err.Error()))
	}
	consumer , err = sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Logger.Error("fail to initialize kafka" , log.Any("error" , err.Error()))
	}
}

func ConsumerMsg (callBack ConsumerCallback){
	partitionConsumer , err := consumer.ConsumePartition(topic , 1 , sarama.OffsetNewest)
	if err != nil {
		log.Logger.Error("iConsumePartition error", log.Any("ConsumePartition error", err.Error()))
		return
	}

	defer partitionConsumer.Close()
	for {
		msg := <-partitionConsumer.Messages()
		if callBack != nil {
			callBack(msg.Value)
		}
	}
}