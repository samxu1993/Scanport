package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

//扫描全部端口


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
	if _, err := conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n")); err != nil {
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
	check_address("github.com:80")
	check_address("58.246.240.66:22222")
}