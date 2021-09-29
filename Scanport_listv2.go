package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func main() {

	var begin =time.Now()
	//wg
	var wg sync.WaitGroup
	//var ip = "192.168.8.1"
	//var ip = "121.37.174.77"
	fi, err := os.Open("ip.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		fmt.Println(string(a))

		file, err := os.OpenFile("port.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
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
				var address = fmt.Sprintf("%s:%d", string(a), i)

				//conn, err := net.DialTimeout("tcp", address, time.Second*10)
				conn, err := net.Dial("tcp", address)
				if err != nil {
					//fmt.Println(address, "是关闭的", err)
					return
				}
				conn.Close()
				fmt.Println("tcp",address, "打开")
				fmt.Println(conn)
				content := []byte("tcp "+address+"\n")
				if _, err = file.Write(content); err != nil {
					fmt.Println(err)
				}
				conn1, err := net.Dial("udp", address)
				if err != nil {
					//fmt.Println(address, "是关闭的", err)
					return
				}
				conn1.Close()
				fmt.Println("udp",address, "打开")
				content1 := []byte("udp "+address+"\n")
				if _, err = file.Write(content1); err != nil {
					fmt.Println(err)
				}


			}(j)
		}
		//等待wg
		wg.Wait()
		var elapseTime = time.Now().Sub(begin)
		fmt.Println("耗时:", elapseTime)



	}
	//fmt.Println("请输入IP ")
	//fmt.Scanln(&ip)
	//循环

}