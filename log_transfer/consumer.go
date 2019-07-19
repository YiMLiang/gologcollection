package main

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"logagent/es"
	"logagent/kafka"
	"time"
)

/**
kafka消费数据逻辑
*/

func run() (err error) {

	//kafka : 取出topic下的所有分区
	partitionList, err := kafka.KClient.Client.Partitions(kafka.KClient.Topic)
	if err != nil {
		logs.Error("log_transfer.consumer : Failed to get the list of kafka partitions: topic = : ,err", kafka.KClient.Topic, err)
		return
	}

	//遍历每个分区下的数据进行消费
	for partition := range partitionList {
		partitionConsumer, e := kafka.KClient.Client.ConsumePartition(kafka.KClient.Topic, int32(partition), sarama.OffsetNewest)
		if e != nil {
			err = e
			logs.Error("log_transfer.consumer : Failed to start consumer for partition %d: %s\n", partition, err)
			return
		}
		//逻辑处理完关掉kafka
		defer partitionConsumer.AsyncClose()

		//针对每个分区下启动一个goroutine消费数据
		go func(partitionConsumer sarama.PartitionConsumer) {
			//每个goroutine中创建一个waitGroup来同步goroutine
			kafka.KClient.Wg.Add(1)
			for msg := range partitionConsumer.Messages() {
				logs.Debug("Partition:%d, Offset:%d, Key:%s, Value:%s", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				err = es.SendToEs(kafka.KClient.Topic, string(msg.Value))
				if err != nil {
					logs.Warn("logagent_transfer.run sendToEs failed, err: %v", err)
				}
			}
			kafka.KClient.Wg.Done()
		}(partitionConsumer)
	}
	time.Sleep(time.Millisecond * 100)
	kafka.KClient.Wg.Wait()
	//err = kafka.KClient.Client.Close()
	//if err != nil {
	//	logs.Error("log_transfer.run : kafka consumer close failed, err : %v", err)
	//	return
	//}
	return
}
