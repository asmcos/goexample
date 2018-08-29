## requests

`requests` 是一个用golang 语言clone python版本的requests库。
golang 自带的net/http功能已经非常完善。它和python里面的urllib系列库一样，功能足够，但是使用起来非常繁琐。

python版本的requests简直就是神器，让很多人非常敬仰。

因此我就动手按照python的requests的思想封装了一个 requests。

动手之前，我想尽量兼容 python的写法。但是golang不支持函数名称小写，
也不支持 参数命名。所以我只能用参数类型的方案来支持动态参数。

## 安装

golang 1.11之前

```
go get -u github.com/asmcos/requests
```

golang 1.11之后，编辑一个go.mod文件

```
module github.com/asmcos/requests
```


## 开始使用

``` go
package main

import (
        "github.com/asmcos/requests"
        "fmt"
)

func main (){

        req := requests.Requests()
        resp := req.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})
        println(resp.Text())
        fmt.Println(resp.R.StatusCode)
        fmt.Println(resp.R.Header["Content-Type"])
}

```

第一个例子为什么是它？ 因为python requests第一个例子是它。。。呵呵

    注意：：： `密码` 和用户名要用github真实用户才能测试。


## 创建请求的方法

### 例子1

极简使用

``` go
package main

import "github.com/asmcos/requests"

func main (){

        resp := requests.Get("http://go.xiulian.net.cn")
        println(resp.Text())
}
```

其实 requests实现都是先用Requests()函数创建一个 request 和 client，
再用req.Get去请求。

requests.Get 是一个封装，对Requests()和req.Get的一个封装。

### 例子2

这个例子是分成2个步骤，来操作，这样的好处是可以通过req来设置各种请求参数。
后面的例子会展示各种设置。

``` go
package main

import "github.com/asmcos/requests"


func main (){

        req := requests.Requests()

        resp := req.Get("http://go.xiulian.net.cn",requests.Header{"Referer":"http://www.jeapedu.com"})

        println(resp.Text())

}
```
