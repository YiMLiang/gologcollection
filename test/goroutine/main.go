package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	w := sync.WaitGroup{}

	w.Add(10)

	for i:=0;i<10;i++{
		work(&w,i)
	}
	w.Wait()
	fmt.Println("success")
}

func work(w *sync.WaitGroup,i int){
	fmt.Println("worker :",i)
	time.Sleep(time.Second)
	w.Done()
}
