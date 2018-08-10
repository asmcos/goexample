## 定义

!!! 解释
     URL Schemes 有两个单词：

     URL，我们都很清楚，http://go.xiulian.net.cn 就是个 URL，我们也叫它链接或网址；

     Schemes，表示的是一个 URL 中的一个位置——最初始的位置，即 ://之前的那段字符。比如 http://www.cpython.org 这个网址的 Schemes 是 http。

```
type URL struct {
    Scheme     string
    Opaque     string    // encoded opaque data
    User       *Userinfo // username and password information
    Host       string    // host or host:port
    Path       string    // path (relative paths may omit leading slash)
    RawPath    string    // encoded path hint (see EscapedPath method)
    ForceQuery bool      // append a query ('?') even if RawQuery is empty
    RawQuery   string    // encoded query values, without '?'
    Fragment   string    // fragment for references, without '#'
}
```

### URL 例子

```
[scheme:][//[userinfo@]host][/]path[?query][#fragment]
```

下面例子，我修改了原始例子，觉得原文档例子有歧义。

```
u := new(url.URL)

u.Scheme = "https"
u.Host = "google.com"
q := u.Query()
q.Set("q", "golang")
u.RawQuery = q.Encode()
fmt.Println(u)

```

输入结果

     https://google.com/search?q=golang


## path

      Path 解码路径
      RawPath 原始路径


```
package main

import (
	"fmt"
	"log"
	"net/url"
)

func main() {
	// Parse + String preserve the original encoding.
	u, err := url.Parse("https://example.com/foo%2fbar")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u.Path)      //解码存储
	fmt.Println(u.RawPath)   //原始存储
	fmt.Println(u.String())
}
```

结果如下：

```
/foo/bar
/foo%2fbar
https://example.com/foo%2fbar
```

## EscapedPath

     func (u *URL) EscapedPath() string

EscapedPath 返回 u.Path 的转义形式。一般来说，任何路径都有多种可能的转义形式。EscapedPath 在 u.Path 有效转义时返回 u.RawPath 。否则，EscapedPath 将忽略 u.RawPath 并自行计算转义表单。 String 和 RequestURI 方法使用 EscapedPath 来构造它们的结果。通常，代码应该调用 EscapedPath ，而不是直接读取 u.RawPath 。


## RequestURI

将path 编码，无论原始URL是否编码都返回编码地址。
下面代码，我写了2个URL的例子。

```
package main

import (
        "fmt"
        "log"
        "net/url"
)

func main() {

        // Parse + String preserve the original encoding.
        u, err := url.Parse("https://example.com/搜索")
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println(u.Path)      //解码存储
        fmt.Println(u.RequestURI())   //压缩存储
        fmt.Println(u.String())

        u, err = url.Parse("https://example.com/%E6%90%9C%E7%B4%A2")
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println(u.Path)      //解码存储
        fmt.Println(u.RequestURI())   //压缩存储
        fmt.Println(u.String())

}

```


结果

```
/搜索
/%E6%90%9C%E7%B4%A2
https://example.com/%E6%90%9C%E7%B4%A2
/搜索
/%E6%90%9C%E7%B4%A2
https://example.com/%E6%90%9C%E7%B4%A2
```


## Values
     type Values map[string][]string

用map来存储query 参数

操作方法如下：

```
v := url.Values{}
v.Set("name", "Ava")
v.Add("friend", "Jess")
v.Add("friend", "Sarah")
v.Add("friend", "Zoe")
// v.Encode() == "name=Ava&friend=Jess&friend=Sarah&friend=Zoe"
fmt.Println(v.Get("name"))
fmt.Println(v.Get("friend"))
fmt.Println(v["friend"])
```

结果。

```
Ava
Jess
[Jess Sarah Zoe]
```

这个功能可以对get参数进行转换编码。
