# requests

`requests` 是一个用golang 语言clone python版本的requests库。
golang 自带的net/http功能已经非常完善。它和python里面的urllib系列库一样，功能足够，但是使用起来非常繁琐。

python版本的requests简直就是神器，让很多人非常敬仰。

因此我就动手按照python的requests的思想封装了一个 requests。

动手之前，我想尽量兼容 python的写法。但是golang不支持函数名称小写，
也不支持 参数命名。所以我只能用参数类型的方案来支持动态参数。

# 安装

golang 1.11之前

```
go get -u github.com/asmcos/requests
```

golang 1.11之后，编辑一个go.mod文件

```
module github.com/asmcos/requests
```

参考[go.mod](https://github.com/asmcos/requests/blob/master/examples/go.mod)例子

# 开始使用

```
req := requests.Requests()
resp := req.Get("https://api.github.com/user",requests.Auth{"asmcos","password...."})
println(resp.Text())
```

第一个例子为什么是它？ 因为python requests第一个例子是他。。。呵呵
