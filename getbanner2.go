package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

//扫描全部端口
func scanport()  {
	var begin =time.Now()
	//wg
	var wg sync.WaitGroup
	var ip string
	//ip = "10.21.31.1"
	//var ip = "121.37.174.77"
	fmt.Println("请输入查询地址或IP ")
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
			//fmt.Println( address, "打开")
			check_address(address)
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

// 假定是 SSH 服务。
// 返回 banner 第一行。
func assume_ssh(address string) (string, error) {
	conn, err := net.DialTimeout("tcp", address, time.Second*10)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	tcpconn := conn.(*net.TCPConn)
	// 设置读取的超时时间
	tcpconn.SetReadDeadline(time.Now().Add(time.Second * 5))
	reader := bufio.NewReader(conn)
	return reader.ReadString('\n')
}

func split_http_head(data []byte, atEOF bool) (advance int, token []byte, err error) {
	head_end := bytes.Index(data, []byte("\r\n\r\n"))
	if head_end == -1 {
		return 0, nil, nil
	}
	return head_end + 4, data[:head_end+4], nil
}

// 假定是 HTTP 服务。
// 返回 "/" HTTP 返回头。
func assume_http(address string) (string, error) {
	conn, err := net.DialTimeout("tcp", address, time.Second*10)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	tcpconn := conn.(*net.TCPConn)
	// 设置写的超时时间
	tcpconn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	if _, err := conn.Write([]byte("HEAD / HTTP/1.1\r\n\r\n")); err != nil {

		return "", err
	}
	// 设置读的超时时间
	tcpconn.SetReadDeadline(time.Now().Add(time.Second * 5))
	scanner := bufio.NewScanner(conn)
	scanner.Split(split_http_head)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	err = scanner.Err()
	if err == nil {
		err = io.EOF
	}
	return "", err
}

func check_address(address string) {
	result := make(chan string, 2)
	done := make(chan int, 1)
	var g sync.WaitGroup
	g.Add(2)
	go func() {
		if r, e := assume_ssh(address); e == nil {
			result <- fmt.Sprintf("SSH: %s", r)
		}
		g.Done()
	}()
	go func() {
		if r, e := assume_http(address); e == nil {
			result <- fmt.Sprintf("HTTP: %s", r)
		}
		g.Done()
	}()
	go func() {
		g.Wait()
		done <- 1
	}()
	select {
	case <-done:
		fmt.Printf("# %s\n无结果\n", address)
	case r := <-result:
		fmt.Printf("# %s\n%s", address, r)
	}
}

func main() {
	scanport()
	//check_address("192.168.1.25:8000")
	//check_address("192.168.1.25:22")
}