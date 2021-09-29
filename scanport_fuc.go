package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func getip(ip string) string {
	return ip
}

func main() {

	var begin =time.Now()
	//wg
	var wg sync.WaitGroup
	var ip string
	//ip = "10.21.31.1"
	//var ip = "121.37.174.77"
	fmt.Println("请输入IP ")
	fmt.Scanln(&ip)
	//循环
	file, err := os.OpenFile(ip+".txt", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for j := 1; j <= 65535; j++ {
		//添加wg
		wg.Add(1)
		go func(i int) {
			//释放wg
			defer wg.Done()
			var address = fmt.Sprintf("%s:%d", ip, i)

			//conn, err := net.DialTimeout("tcp", address, time.Second*10)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				//fmt.Println(address, "是关闭的", err)
				return
			}
			conn.Close()
			fmt.Println( address, "打开")

			content := []byte(address+"\n")
			if _, err = file.Write(content); err != nil {
				fmt.Println(err)
			}
		}(j)
	}
	//等待wg
	wg.Wait()
	var elapseTime = time.Now().Sub(begin)
	fmt.Println("耗时:", elapseTime)
}