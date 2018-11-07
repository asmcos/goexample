## gin

gin 是一款轻量级的 golang web server框架。

源代码地址 [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)

目前在github上star最多。这在一定层面上表示认可度最高。

## 开始

   go get github.com/gin-gonic/gin

## 写一个example.go

```go
package main

import (
        "net/http"

        "github.com/gin-gonic/gin"
)

func main (){

        r := gin.Default()

        r.GET("/hello",func(c *gin.Context){
                c.String(http.StatusOK,"Hello gocoin!")
        })

        r.Run(":8080")
}

```

执行 go run example.go

打开浏览器访问： http://127.0.0.1:8080/hello

## 路由

## 模版

## 数据库

## cookies

## 验证用户
