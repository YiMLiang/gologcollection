package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Result struct {
	resp *http.Response
	err error
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	tr := &http.Transport{}
	client := &http.Client{Transport:tr}

	request, e := http.NewRequest("GET", "https://www.baidu.com", nil)
	if e!=nil {
		fmt.Println("http request failed err:",e)
		return
	}

	//建个管道
	c:=make(chan Result,1)

	//goroutine 中不一定非要做http请求 也可以做别的事情，比如：任何任务都可以用ctx控制超时
	go func() {
		time.Sleep(time.Second*3)
		resp, err := client.Do(request)
		result := Result{resp:resp,err:err}
		c<-result

	}()

	select {
	case <-ctx.Done():
		tr.CancelRequest(request)
		x := <-c

		fmt.Println("time out!",x)

	case Res:=<-c:
		defer Res.resp.Body.Close()
		bytes, _ := ioutil.ReadAll(Res.resp.Body)
		fmt.Printf("server response is : %s",bytes)
	}

	return
}