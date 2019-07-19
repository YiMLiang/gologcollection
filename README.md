# gologcollection
Golang log collection【Go实现日志tail】

### 初衷是想达到一个日志收集-读日志-发送到es这样一个过程

1.日志收集的路径初步设定死是在某个路径下，如果想修改收集的日志还需停了程序修改后再起
2.添加etcd模块，目的就是类似于zookeeper吧，做一个服务发现和注册的功能，针对某台机器，设定一个对应的key，将要收集的日志路径写到对应的value中
  再起一个watch去监听每个key下面的配置变化，若有改变则监听到之后修改，这样不必手动调整
3.添加消息中间件kafka，因为直接读日志发到es有点单一，而且没有一个中间件做负载也不太安全，可能会造成数据丢失，所以中间用kafka做了一个消息中介，
  读到日志后实时发送给kafka，kafka拿到数据给es消费
4.kibana做es的web界面展示，提供了便捷的关键字查询功能，很方便
5.写了web界面用来配置项目信息和日志的信息【日志信息即etcd设置的value 即要收集的日志的路径】


### flow chart

![](ER.png)

### technology
Golang,etcd,kafka,ElasticSearch,kibana




