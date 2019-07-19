package test

import (
	"context"
	"fmt"
)

func process (ctx context.Context){
	ret,ok := ctx.Value("trace_id").(int)
	if !ok {
		ret = 123123
	}
	fmt.Println("ret = :",ret)
	str,_ := ctx.Value("session_id").(string)

	fmt.Printf("ret = : %v,session = : %v",ret,str)
}

func main() {

	ctx := context.WithValue(context.Background(), "trace_id", 22222)

	ctx = context.WithValue(ctx,"session_id","lff")

	process(ctx)
}



