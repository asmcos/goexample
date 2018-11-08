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
数据库模块常用的是 gorm。

安装 gorm
    go get github.com/jinzhu/gorm
    go get github.com/mattn/go-sqlite3

因为我例子用的sqlite3，所以也安装了sqlite3 依赖库。

       bogon:gocoin jiashenghe$ cat gorm_example.go

```go

package main

import (
        "github.com/jinzhu/gorm"
        _ "github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
)

type Person struct {
        gorm.Model
        FirstName string
        LastName  string
}

func main() {
        db, _ := gorm.Open("sqlite3", "./gorm.db")
        defer db.Close()

        db.AutoMigrate(&Person{})

        p1 := Person{FirstName: "John", LastName: "Doe"}
        p2 := Person{FirstName: "Jane", LastName: "Smith"}

        db.Create(&p1)
        var ps []Person

	      fmt.Println("-------1------")
        db.Find(&ps)
      	for _,k := range ps{
      		fmt.Println(k)
      	}


        db.Create(&p2)
      	fmt.Println("-------2------")
        db.Find(&ps)
      	for _,k := range ps{
      		fmt.Println(k)
      	}
}

```

执行结果：

```
-------1------
{{1 2018-11-08 10:12:48.913042 +0800 +0800 2018-11-08 10:12:48.913042 +0800 +0800 <nil>} John Doe}
-------2------
{{1 2018-11-08 10:12:48.913042 +0800 +0800 2018-11-08 10:12:48.913042 +0800 +0800 <nil>} John Doe}
{{2 2018-11-08 10:12:48.914945 +0800 +0800 2018-11-08 10:12:48.914945 +0800 +0800 <nil>} Jane Smith}
```


# gorm + gin
合并代码如下：
```go
package main

import (
        _ "net/http"

        "github.com/gin-gonic/gin"
        "github.com/jinzhu/gorm"
        _ "github.com/jinzhu/gorm/dialects/sqlite"
        "fmt"
)

var db *gorm.DB
var err error

type Person struct {
        gorm.Model
        FirstName string
        LastName  string
}

func init_db (){

        db, err = gorm.Open("sqlite3", "./gorm.db")
        if err != nil {
                fmt.Println(err)
        }
        db.AutoMigrate(&Person{})

}

func main (){

        init_db()
        defer db.Close()

        r := gin.Default()

        r.GET("/hello",func(c *gin.Context){
                //c.String(http.StatusOK,"Hello gocoin!")
                var people []Person
                if err := db.Find(&people).Error; err != nil {
                        c.AbortWithStatus(404)
                        fmt.Println(err)
                } else {
                        c.JSON(200, people)
                }


        })

        r.Run(":8080")
}
```

go run main.go

在浏览器浏览 http://127.0.0.1:8080/hello

结果如下:

```
[{"ID":1,"CreatedAt":"2018-11-08T10:12:48.913042+08:00","UpdatedAt":"2018-11-08T10:12:48.913042+08:00","DeletedAt":null,"FirstName":"John","LastName":"Doe"},{"ID":2,"CreatedAt":"2018-11-08T10:12:48.914945+08:00","UpdatedAt":"2018-11-08T10:12:48.914945+08:00","DeletedAt":null,"FirstName":"Jane","LastName":"Smith"}]
```
## cookies

## 验证用户
