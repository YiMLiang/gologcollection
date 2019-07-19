package config

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
)

//k8吹一遍，一致性数据etcd吹一遍，gin，beego，xorm，gomicro grpc都吹一遍
var (
	AppConfig *Config
)

type Config struct {
	LogLevel    string
	LogPath     string
	KafkaAddr   string
	Topic       string //加载kafka topic
	EtcdAddr    string //加载etcd addr
	EtcdKey     string //加载etcd key
	EsAddr      string //加载es addr
	ChanSize    int
	CollectConf [] CollectConf
}
type CollectConf struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

/**
加载配置文件
*/
func LoadConf(configType, fileName string) (err error) {

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

	AppConfig.ChanSize, err = conf.Int("collect::chan_size")
	if err != nil {
		AppConfig.ChanSize = 100
	}

	//加载kafka配置

	AppConfig.KafkaAddr = conf.String("kafka::server_addr")
	if len(AppConfig.KafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka addr")
		return
	}

	//加载etcd配置
	AppConfig.EtcdAddr = conf.String("etcd::etcd_addr")
	if len(AppConfig.EtcdAddr) == 0 {
		err = fmt.Errorf("invalid etcd_addr addr")
		return
	}

	AppConfig.EtcdKey = conf.String("etcd::etcd_key")
	if len(AppConfig.EtcdKey) == 0 {
		err = fmt.Errorf("invalid etcd_key addr")
		return
	}

	//日志收集相关
	err = loadCollectConf(conf)
	if err != nil {
		fmt.Printf("load collect conf failed , err: %v\n", err)
	}
	return
}

/**
业务日志收集
*/
func loadCollectConf(conf config.Configer) (err error) {
	var cc CollectConf
	cc.LogPath = conf.String("collect::log_path")
	fmt.Println(cc.LogPath)
	if len(cc.LogPath) == 0 {
		err = errors.New("invalid collect :: log_path")
		return
	}
	cc.Topic = conf.String("collect::topic")
	fmt.Println(cc.Topic)
	if len(cc.Topic) == 0 {
		err = errors.New("invalid collect :: topic")
		return
	}

	AppConfig.CollectConf = append(AppConfig.CollectConf, cc)
	return
}
