package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379","localhost:22379"},
		DialTimeout: 5 * time.Second,
	})

	if err!=nil {
		fmt.Println("connetc failed,err:",err)
		return
	}

	fmt.Println("connect success")
	defer client.Close()

	for {
		rch := client.Watch(context.Background(),"/logagent/conf/")
		for wresp := range rch{
			for _,event := range wresp.Events{
				fmt.Printf("%s %q : %q\n",event.Type,event.Kv.Key,event.Kv.Value)
			}
		}
	}
	
}
