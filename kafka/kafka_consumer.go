package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"strings"
	"sync"
)

var (
	KClient *KafkaClient
)

type KafkaClient struct {
	Client sarama.Consumer
	Addr   string
	Topic  string
	Wg     sync.WaitGroup
}
/**
加载kafkaConsumer
 */
func InitKafkaConsumer(addr string, topic string) (err error) {
	KClient = &KafkaClient{}
	consumer, err := sarama.NewConsumer(strings.Split(addr, ","), nil)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}
	KClient.Client = consumer
	KClient.Addr = addr
	KClient.Topic = topic
	return
}
