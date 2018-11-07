package main

import (
    "log"
    "github.com/asmcos/tlsdump"
)

func main() {
    log.SetFlags(log.Lshortfile)

    conf := &tlsdump.Config{
         InsecureSkipVerify: true,
    }

    conn, err := tlsdump.Dial("tcp", "127.0.0.1:443", conf)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    n, err := conn.Write([]byte("hello\n"))
    if err != nil {
        log.Println(n, err)
        return
    }

    buf := make([]byte, 100)
    n, err = conn.Read(buf)
    if err != nil {
        log.Println(n, err)
        return
    }

    println(string(buf[:n]))
}

