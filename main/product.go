package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logagent/kafka"
	"logagent/tailf"
	"time"
)

/**
循环获取日志文件中的每一行
*/
func serverRun() (err error) {
	for {
		msg := tailf.GetOneLine() //得到每一行日志信息
		err = sendToKafka(msg)    //发送kafka
		if err != nil {
			logs.Error("send to Kafka failed,err : %v", err)
			time.Sleep(time.Millisecond * 10)
			continue
		}
	}
	return
}

func sendToKafka(msg *tailf.TextMsg) (err error) {
	fmt.Printf("sendToKafka : read msg : %s", msg.Msg)
	err = kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
