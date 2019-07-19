package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"github.com/tsuna/gohbase/pb"
	"os"
)

func init() {

	// 以Stdout为输出，代替默认的stderr
	logrus.SetOutput(os.Stdout)
	// 设置日志等级
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {

	client := gohbase.NewClient("192.168.150.134")

	getRequest, _ := hrpc.NewGetStr(context.Background(), "emp", "1")
	getRsp, _ := client.Get(getRequest)

	for _, cell := range getRsp.Cells{
		fmt.Println(string((*pb.Cell)(cell).GetFamily()))
		fmt.Println(string((*pb.Cell)(cell).GetQualifier()))
		fmt.Println(string((*pb.Cell)(cell).GetValue()))
		fmt.Println((*pb.Cell)(cell).GetTimestamp())
	}

}
