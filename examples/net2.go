package main

import (
	"net"
	"fmt"
)

func handleConnection(conn net.Conn) {  
	buf := make([]byte,2000)

	n,_ := conn.Read(buf)
	fmt.Println(conn.RemoteAddr().String())
	fmt.Println(n)
}

func main (){
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		fmt.Println(conn)
		go handleConnection(conn)
	}

}

