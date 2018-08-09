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


## Listen 创建一个服务

```
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
```

这个处理和其他语言的网络没有太多区别，就是Listen 8080端口，等待链接。
当Accept一个链接后就是 go(Goroutines)一个线程去处理这个链接conn。


## 域名解析（名字）

无论我们是调用Dial 来获得解析，还是直接调用LookupHost和LookupAddr来解析，这都是和OS息息相关的。

在Unix系列的OS上，解析器有两个方法解析域名。方法一就是用go发送到dns解析（dns在/etc/resolv.conf文件中配置）。方法二是cgo的解析器调用C库例程，如getaddrinfo和getnameinfo。

调用方法一时如果出现问题，会消耗 Goroutines 线程。而调用方法二如果出现出问题会消耗操作系统线程。因此go默认的情况是使用方法一。


可以通过将GODEBUG环境变量的netdns值（请参阅包运行时）设置为go或cgo来覆盖解析器决策，如下所示：

```
export GODEBUG = netdns = go＃force pure Go resolver
export GODEBUG = netdns = cgo＃force cgo resolver
```

在Windows上，解析器始终使用C库函数，例如GetAddrInfo和DnsQuery。

## func JoinHostPort

     func JoinHostPort(host, port string) string

JoinHostPort 函数的目的是将host和port合成一个“host:port”格式的网络地址，如果IPV6这种的自身带有“：”，返回的格式应该是 "[host]:port"

## LookupAddr

       func LookupAddr(addr string) (names []string, err error)

LookupAddr 函数是给定一个 地址反向查找，返回是的改名字对应的映射列表。

```
lt, _ := net.LookupAddr("127.0.0.1")
fmt.Println(lt) //[localhost],根据地址查找到改地址的一个映射列表

```
