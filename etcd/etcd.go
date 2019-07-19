package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"logagent/config"
	"logagent/ip"
	"logagent/tailf"
	"strings"
	"time"
)

/**
go连接etcd
1.初始化etcd
2.获取 etcd中设置的配置
3.监听etcd 中配置的变化 并启与keys数量相同的goroutine来监听key的变化
连上etcd 之后将配置文件读出来
配置时： key 要带上ip，即 机器+ip 才能对应这台机器上面的配置文件
*/

type EtcdClient struct {
	client *clientv3.Client //client类型
	keys   []string         // etcd 的 key
}

var (
	etcdClient *EtcdClient
)
/**
初始化etcd
*/

func InitEtcd(addr string, key string) (collectConf []config.CollectConf, err error) {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("etcd.InitEtcd : connetc etcd failed,err:", err)
		return
	}
	fmt.Println("etcd.InitEtcd : connect etcd success!!!")

	//设置全局etcdClient
	etcdClient = &EtcdClient{
		client: client,
	}

	collectConf, err = getConfig(client, addr, key)
	return collectConf, err
}

/**
获取 etcd中设置的配置
*/

func getConfig(client *clientv3.Client, addr string, key string) (collectConf []config.CollectConf, err error) {
	//判断配置文件是否是 以 "/" 结尾，若不是则在后面加上 "/"
	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}

	for _, ipstr := range ip.LocalIpArray {
		//拼一下每一台机器对应的 etcd_key
		etcdKey := fmt.Sprintf("%s%s", key, ipstr)
		//将拼好的key 加到全局的etcdClient.keys 数组中
		etcdClient.keys = append(etcdClient.keys, etcdKey)
		//设置一个连接超时的 ctx
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		//获取 etcd 配置
		resp, err := client.Get(ctx, etcdKey)
		if err != nil {
			logs.Error("etcd.getConfig : get conf from etcd failed,err :%v", err)
			continue
		}
		cancel()

		logs.Debug("etcd.getConfig : success ,resp from etcd :%v", resp.Kvs)

		//循环获取从etcd中 get到对应conf的 ，每一个 key-value配置 ｛conf可能是一个，可能是一个数组-包含多个｝
		for _, v := range resp.Kvs {
			//fmt.Printf("key : %s, value : %s", k, v)
			//判断从etcd中取出的key是否和你拼成的key相同
			if string(v.Key) == etcdKey {
				//拿到配置，取它的值(value)，封装成json格式 此时 key不变，value变为了json格式
				err := json.Unmarshal(v.Value, &collectConf)
				if err != nil {
					logs.Error("etcd.getConfig : json.Unmarshal failed err:%v", err)
					continue
				}

				logs.Debug("etcd.getConfig : 从etcd中取出对应 key : %s\n 配置 config is : %s", v.Key, collectConf)
			}

		}
	}

	//监听etcd 中配置的变化
	initEtcdWatcher()

	return
}
/**
加载etcd监听器
 */
func initEtcdWatcher() {
	//getConfig 方法中已经将所有拼成的 etcdKey 放[append]到 etcdClient.keys数组中，遍历key并监听其变化
	//达到实时获取配置的功能
	for _, key := range etcdClient.keys {
		logs.Debug("etcd.initEtcdWatcher : key = %v:", key)
		go watchKey(key)
	}
}

/**
etcd监听器实现
 */
func watchKey(key string) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.AppConfig.EtcdKey},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("etcd.watchKey : connetc failed,err:", err)
		return
	}

	logs.Debug("etcd.watchKey : 开始监测etcd节点,key = :%s", key)

	for {
		rch := client.Watch(context.Background(), key)
		var collectConf []config.CollectConf
		var getConfSuccess = true
		//遍历 WatchChan 管道,从管道里面取出所监听的key
		//接着遍历key所对应的每个conf
		for wresp := range rch {
			for _, event := range wresp.Events {
				//fmt.Printf("type : %s ,key : %q, value : %q\n", event.Type, event.Kv.Key, event.Kv.Value)
				if event.Type == mvccpb.DELETE {
					logs.Warn("etcd.watchKey : key[%s] 's config deleted", key)
					continue
				} else if event.Type == mvccpb.PUT && string(event.Kv.Key) == key {
					//key相同又为put时 将新配置设置进etcd中
					err = json.Unmarshal(event.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key [%s], Unmarshal[%s], err:%v ", key, event.Kv.Value, err)
						getConfSuccess = false
						continue
					}
				} else {
					logs.Debug("etcd.watchKey 暂无匹配项,请检查")
				}

				logs.Debug("etcd.watchKey : get type and config from etcd, type = [%s], kv = : %q , %q\n", event.Type, event.Kv.Key, event.Kv.Value)

			}

			if getConfSuccess {
				logs.Debug("etcd.watchKey : get config from etcd succ, %v", collectConf)
				//得到新配置后更新原有配置
				tailf.UpdateConfig(collectConf)
			}
		}
	}
}
