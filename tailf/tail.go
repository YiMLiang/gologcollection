package tailf

import (
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"logagent/config"
	"sync"
	"time"
)

type TailObj struct {
	tail     *tail.Tail         //tail对象
	conf     config.CollectConf //日志结构体 包含 logpath 和topic
	status   int                //状态
	exitChan chan int           //若logpath被删除或被替换，exitChan作为标识存在,结束的标志
}

const (
	StatusNormal = 1 //状态正常表示logpath没有发生改变
	StatusDelete = 2 //delete表示logpath已被删除
)

/**
Msg：实际发送的msg，
topic：日志写到哪个topic里面
*/
type TextMsg struct {
	Msg   string
	Topic string
}

type TailObjManager struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg
	lock     sync.Mutex
}

var (
	tailObjMgr *TailObjManager
)

/**
加载tail组件
*/

func InitTail(conf []config.CollectConf, chanSize int) (err error) {
	tailObjMgr = &TailObjManager{
		//假设写死100条
		msgChan: make(chan *TextMsg, chanSize),
	}

	if len(conf) == 0 {
		logs.Error("tail.InitTail : invalid conf for collect ,err: %v", err)
		return
	}

	for _, v := range conf {
		createNewTask(v)
	}

	return
}

/**
从tail中逐行读取日志文件
*/
func readFromTail(tailObj *TailObj) {

	for true {
		select {
		//如果读到tail的一行
		case line, ok := <-tailObj.tail.Lines:
			if !ok {
				logs.Warn("tail.go.readFromTail : tail file close reopen,fileName : %s", tailObj.tail.Filename)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			//1.组织发送信息结构体textMsg
			textMsg := &TextMsg{
				Msg:   line.Text,
				Topic: tailObj.conf.Topic,
			}
			//2.将信息结构体textMsg放到chan中
			tailObjMgr.msgChan <- textMsg
		//如果只读到exitChan中的 1
		case <-tailObj.exitChan:
			logs.Warn("tail obj will exited ,conf :%v", tailObj.conf)
			return
		}
	}
}

/**
从  channel 中获取一行数据
*/
func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return msg
}

/**
更新新修改的配置文件
*/
func UpdateConfig(config []config.CollectConf) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()

	for _, conf := range config {
		var flag = false
		for _, elem := range tailObjMgr.tailObjs {
			if conf.LogPath == elem.conf.LogPath {
				//并无变化
				flag = true
				break
			}
		}
		//若配置无变化:continue
		if flag {
			continue
		}
		//当配置发生变化时 创建一次新的tail 继续读新配置中设定的日志文件
		createNewTask(conf)
	}

	//新建tailObjs对象装新的tail信息，可能是未变化的tail，即配置文件没变，
	//可能是新的tail，即配置变了，可能要从其他路径读日志了
	var tailObjs []*TailObj
	//遍历tailObjs对象 添加是否退出标识 exitChan 若退出则向chan中放1 默认status为2，若不退出则设置status为1
	for _, obj := range tailObjMgr.tailObjs {
		obj.status = StatusDelete
		for _, oneConf := range config {
			if oneConf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}
		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue
		}
		//将修改的加到新 tailObjs 数组中.
		tailObjs = append(tailObjs, obj)
	}
	//修改后将修改后的tailObjs 覆盖之前的 tailObjs
	tailObjMgr.tailObjs = tailObjs
	return
}

/**
创建新的tail，读配置中设定的日志文件
*/
func createNewTask(conf config.CollectConf) {
	obj := &TailObj{
		conf:     conf,
		exitChan: make(chan int, 1),
	}
	tails, errTail := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Poll:      true,
	})

	if errTail != nil {
		logs.Error("collect fileName[%s] failed ,err :%v", conf.LogPath, errTail)
		return
	}

	obj.tail = tails
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

	//启动一个goroutine 读日志
	go readFromTail(obj)
}
