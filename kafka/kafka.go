package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"time"
)

/**
加载kafka配置
 */

var (
	producer sarama.SyncProducer
)

func InitKafka(kafkaAddr string) (e error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	producer, e = sarama.NewSyncProducer([]string{kafkaAddr,"192.168.150.134:9092","192.168.150.136:9092"}, config)
	if e != nil {
		logs.Error("init kafka producer failed ,err:", e)
		return
	}

	logs.Info("init kafka producer success ")

	return
}

/**
tail日志，逐行读到日志发送kafka逻辑实现
data : 发送kafka数据体
topic : kafka topic
 */
func SendToKafka(data,topic string )(err error){
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	_, _, err = producer.SendMessage(msg)
	if err!=nil{
		logs.Error("send msg failed,err : ,data,topic",err,data,topic)
	}
	//logs.Debug("pid : %v,offset : %v , topic : %v",pid,offset,topic)
	time.Sleep(time.Millisecond*10)
	return
}