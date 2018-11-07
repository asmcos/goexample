package main

import (
    "log"
    "github.com/asmcos/tlsdump"
    "net"
    "bufio"
)

func main() {
    log.SetFlags(log.Lshortfile)

    cer, err := tlsdump.LoadX509KeyPair("server.crt", "server.key")
    if err != nil {
        log.Println(err)
        return
    }

    config := &tlsdump.Config{Certificates: []tlsdump.Certificate{cer}}
    ln, err := tlsdump.Listen("tcp", ":443", config) 
    if err != nil {
        log.Println(err)
        return
    }
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Println(err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    r := bufio.NewReader(conn)
    for {
        msg, err := r.ReadString('\n')
        if err != nil {
            log.Println(err)
            return
        }

        println(msg)

        n, err := conn.Write([]byte("world\n"))
        if err != nil {
            log.Println(n, err)
            return
        }
    }
}
