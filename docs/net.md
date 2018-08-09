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

## 寻找主机名

    LookupCNAME 返回规范名，符合DNS的CNAME规则的。
    LookupHost或LookupIP 返回非规范名

##  SplitHostPort

     func SplitHostPort(hostport string) (host, port string, err error)

把网络格式分析成host 和port ，网络格式包含："host:port", "host%zone:port", "[host]:port" or "[host%zone]:port"


## Buffers

    type Buffers [][]byte
    func (v *Buffers) Read(p []byte) (n int, err error)
    func (v *Buffers) WriteTo(w io.Writer) (n int64, err error)

包含0个，1个或者多个 要写的内容。 writev，read的时候使用。


## Conn


```
type Conn interface {
    // Read reads data from the connection.
    // Read can be made to time out and return an Error with Timeout() == true
    // after a fixed time limit; see SetDeadline and SetReadDeadline.
    Read(b []byte) (n int, err error) //从连接读书句

    // Write writes data to the connection.
    // Write can be made to time out and return an Error with Timeout() == true
    // after a fixed time limit; see SetDeadline and SetWriteDeadline.
    Write(b []byte) (n int, err error) //发送数据到连接

    // Close closes the connection.
    // Any blocked Read or Write operations will be unblocked and return errors.
    Close() error

    // LocalAddr returns the local network address.
    LocalAddr() Addr  //本机地址

    // RemoteAddr returns the remote network address.
    RemoteAddr() Addr //远程地址

    // SetDeadline sets the read and write deadlines associated
    // with the connection. It is equivalent to calling both
    // SetReadDeadline and SetWriteDeadline.
    //
    // A deadline is an absolute time after which I/O operations
    // fail with a timeout (see type Error) instead of
    // blocking. The deadline applies to all future and pending
    // I/O, not just the immediately following call to Read or
    // Write. After a deadline has been exceeded, the connection
    // can be refreshed by setting a deadline in the future.
    //
    // An idle timeout can be implemented by repeatedly extending
    // the deadline after successful Read or Write calls.
    //
    // A zero value for t means I/O operations will not time out.
    SetDeadline(t time.Time) error  //操作的timeout值，0表示不timeout

    // SetReadDeadline sets the deadline for future Read calls
    // and any currently-blocked Read call.
    // A zero value for t means Read will not time out.
    SetReadDeadline(t time.Time) error

    // SetWriteDeadline sets the deadline for future Write calls
    // and any currently-blocked Write call.
    // Even if write times out, it may return n > 0, indicating that
    // some of the data was successfully written.
    // A zero value for t means Write will not time out.
    SetWriteDeadline(t time.Time) error
}

```

## DialTimeout

带Timeout 的Dial

##  Pipe
    func Pipe（）（Conn，Conn）

在内存里将两个conn 同步，建立一个全双工网络。读的一端和写的一端是完全匹配的。两者的数据直接复制，没有任何缓存。

例如： Pipe(Conn1,Conn2)  
如果程序向Conn1 写的内容，通过Conn2能读取。
如果程序向Conn2 写的内容，通过Conn1也能读取。
