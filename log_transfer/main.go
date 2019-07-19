package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logagent/config"
	"logagent/es"
	"logagent/kafka"
	"logagent/log"
)

/**
消费kafka数据 ——> 写入es
*/

func main() {

	fileName := "./conf/logagent_transfer.conf"

	//加载logagent_transfer.conf 日志文件
	err := config.LoadConfTransfer("ini", fileName)
	if err != nil {
		fmt.Println("logaget_transfer : load conf error")
		panic("logaget_transfer : load conf failed")
		return
	}

	//加载日志
	err = log.InitLogger(config.AppConfig.LogPath, config.AppConfig.LogLevel)
	if err != nil {
		fmt.Println("logaget_transfer : init logger error")
		panic("logaget_transfer : init logger failed")
		return
	}

	//初始化kafka
	err = kafka.InitKafkaConsumer(config.AppConfig.KafkaAddr, config.AppConfig.Topic)
	if err != nil {
		fmt.Println("logaget_transfer : init kafka error")
		panic("logaget_transfer : init kafka failed")
		return
	}
	fmt.Println("init kafka success")

	//初始化es
	err = es.InitEs(config.AppConfig.EsAddr)
	if err != nil {
		fmt.Println("logaget_transfer : init es error")
		panic("logaget_transfer : init es failed")
		return
	}
	fmt.Println("init es success")

	logs.Debug("log_transfer : Init All Success!!")

	//启动kafka消费者消费数据逻辑
	err = run()
	if err != nil {
		logs.Error("log_transfer : project run failed, err : %v", err)
		return
	}

	logs.Warn("log_transfer : warning, log_transfer is exited")
}
