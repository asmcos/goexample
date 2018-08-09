package main

import (
	"net"
	"fmt"
	"time"
	// "bufio"
)

func main (){

	time.Sleep(10*time.Second)
	_, err := net.Dial("tcp", "go.xiulian.net.cn:80")
	if err != nil {
		// handle error
	}

	
	_, err1 := net.Dial("tcp", "go.xiulian.net.cn:80")
	if err1 != nil {
		fmt.Println( "handle error" )
	}

	//fmt.Println(conn)
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//status, _ := bufio.NewReader(conn).ReadString('\n')

	//fmt.Println(status)
	return 
}

