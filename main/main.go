package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"logagent/config"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/log"
	"logagent/tailf"
)

func main() {
	//初始化日志相关

	fileName := "./conf/logagent.conf"

	err := config.LoadConf("ini", fileName)
	if err != nil {
		fmt.Println("load conf error")
		panic("load conf failed")
		return
	}

	//加载日志
	err = log.InitLogger(config.AppConfig.LogPath, config.AppConfig.LogLevel)
	if err != nil {
		fmt.Println("init logger error")
		panic("init logger failed")
		return
	}

	//初始化etcd

	collectConf, err := etcd.InitEtcd(config.AppConfig.EtcdAddr, config.AppConfig.EtcdKey)
	if err != nil {
		fmt.Println("init etcd error")
		panic("init etcd failed")
		return
	}
	logs.Debug("Init etcd Success!!")

	err = tailf.InitTail(collectConf, config.AppConfig.ChanSize)
	if err != nil {
		fmt.Println("init tail error")
		panic("init tail failed")
		return
	}

	//初始化kafka
	err = kafka.InitKafka(config.AppConfig.KafkaAddr)
	if err != nil {
		fmt.Println("init kafka error")
		panic("init kafka failed")
		return
	}

	logs.Debug("Init All Success!!")

	//go func() {
	//	var count int
	//	for {
	//		count++
	//		logs.Info("test for logger %d",count)
	//		time.Sleep(time.Second)
	//	}
	//}()

	err = serverRun()
	if err != nil {
		logs.Error("serverRun faild,err:%v", err)
		return
	}

}
