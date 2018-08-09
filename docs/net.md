## 引入包

```
import "net"
```

net 包是提供了一个可以移植的 network I/O 接口。包含 TCP/IP，UDP，域名解析，和Unix domain
sockets 编程。 （注：这些在C语言和其他编程语言都有）

虽然 "net" 包提供了最底层的网络访问接口，但是大部clients编程需要的是基本接口例如：Dial，Listen，Accept和相关的Conn,Listener接口。另外"crypto/tls" 包也有类似：Dial，Listen函数。

## Dial 例子

```
package main

import (
        "net"
        "fmt"
        "bufio"
)

func main (){
        // 原文档golang.org 在国内访问不了，我换了一个域名
        conn, err := net.Dial("tcp", "go.xiulian.net.cn:80")
        if err != nil {
                // handle error
        }
        fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
        status, _ := bufio.NewReader(conn).ReadString('\n')

        fmt.Println(status)
        return
}
```

!!! Dial源代码都调用了哪些系统调用
    sock_posix.go  文件里有posix 的系统调用方法

    * 创建socket
    * connect
    * getsockname
    * getpeername
