package es

import (
	"fmt"
	"gopkg.in/olivere/elastic.v2"
)

/**
go连接es
1.加载es配置
 */

/**
发送es字段定义
 */
type LogMessage struct {
	//App     string
	Topic   string	//topic
	Message string	//发送es消息体
}

var (
	EsClient *elastic.Client
)

/**
加载es相应配置
 */
func InitEs(esAddr string) (err error) {

	client, err := elastic.NewClient(
		elastic.SetSniff(false), elastic.SetURL(esAddr))
	if err != nil {
		fmt.Println("es.InitEs : connect es error, err:", err)
		return
	}
	EsClient = client
	fmt.Println("es.InitEs : connect es success")
	return
}

/**
kafka 消费数据发送到 es 处理逻辑
 */
func SendToEs(topic string, data string) (err error){
	msg := LogMessage{}
	msg.Topic = topic
	msg.Message = data
	_, err = EsClient.Index().
		Index(topic).
		Type(topic).
		Id().
		//Id().
		BodyJson(msg).
		Do()
	if err != nil {
		// Handle error
		panic(err)
		return
	}
	return
}
