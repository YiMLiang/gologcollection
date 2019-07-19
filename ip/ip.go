package ip

import (
	"fmt"
	"net"
)

/**
拿到本机ip 放到 LocalIpArray 数组中
*/

var (
	LocalIpArray []string
)

func init() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get localIp failed,err :", err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			LocalIpArray = append(LocalIpArray, ipnet.IP.String())
		}
	}
	fmt.Println(LocalIpArray)
}
