package config

import (
	"fmt"
	"github.com/astaxie/beego/config"
)

//k8吹一遍，一致性数据etcd吹一遍，gin，beego，xorm，gomicro grpc都吹一遍
var (
	TransferConfig *Config
)

/**
加载transfer所需配置文件
*/
func LoadConfTransfer(configType, fileName string) (err error) {

	conf, err := config.NewConfig(configType, fileName)
	if err != nil {
		fmt.Println("new conf faild , err=", err)
		return
	}

	AppConfig = &Config{}

	//配置日志相关 logs
	AppConfig.LogLevel = conf.String("logs::log_level")
	if len(AppConfig.LogLevel) == 0 {
		AppConfig.LogLevel = "debug"
	}

	AppConfig.LogPath = conf.String("logs::log_path")
	if len(AppConfig.LogPath) == 0 {
		AppConfig.LogPath = "./logs"
	}

	//加载kafka配置

	AppConfig.KafkaAddr = conf.String("kafka::server_addr")
	if len(AppConfig.KafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka addr")
		return
	}

	AppConfig.Topic = conf.String("kafka::topic")
	if len(AppConfig.Topic) == 0 {
		err = fmt.Errorf("invalid topic addr err,please check topic")
		return
	}

	//加载es配置
	AppConfig.EsAddr = conf.String("es::es_addr")
	if len(AppConfig.EsAddr) == 0 {
		err = fmt.Errorf("invalid es addr")
		return
	}

	return
}
