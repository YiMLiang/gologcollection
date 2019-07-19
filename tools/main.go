package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"logagent/config"
	"time"
)

const (
	EtcdKey = "/logagent/conf/10.254.124.191" //本地
	//EtcdKey = "/ymliang/work/logagent/config/10.254.124.191" //本地

	//EtcdKey = "/logagent/conf/192.168.150.134" //虚拟机 master
	//EtcdKey = "/logagent/conf/192.168.162.1"   //虚拟机 master

)

//type LogConf struct {
//	Path  string `json:"path"`
//	Topic string `json:"topic"`
//}

func setLogConfToEtcd() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("setLogConfToEtcd : connetc failed,err:", err)
		return
	}

	fmt.Println("setLogConfToEtcd : connect success")
	defer client.Close()

	var logConfArr []config.CollectConf

	logConfArr = append(
		logConfArr,
		config.CollectConf{
			LogPath: "D:/goprojects/src/logagent/logs/logagent.log",
			//LogPath: "F:/nginx/nginx-1.8.0/logs/access.log",
			Topic:   "nginx_log",
		},
	)

	logConfArr = append(
		logConfArr,
		config.CollectConf{
			LogPath: "D:/project/nginx/logs/error.log",
			Topic:   "nginx_log_err",
		},
	)

	data, err := json.Marshal(logConfArr)
	if err != nil {
		fmt.Println("setLogConfToEtcd : json transform failed!")
	}

	//因为生成etcd 有超时的控制，所以要生成一个ctx做超时处理
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	//操作etcd
	//_, err = client.Put(ctx, "/logagent/conf/", "sample_value")
	_, err = client.Put(ctx, EtcdKey, string(data))
	//操作完成 取消ctx
	cancel()
	if err != nil {
		fmt.Println("clientv3 put failed, err :", err)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*60)
	resp, err := client.Get(ctx, EtcdKey) //withPrefix()是未了获取该key为前缀的所有key-value

	cancel()
	if err != nil {
		fmt.Println("clientv3 get failed ,err :", err)
		return
	}

	for _, v := range resp.Kvs {
		fmt.Printf("key = :%s, value = %s", v.Key, v.Value)
	}
}

func main() {
	setLogConfToEtcd()
}
