## http

http 包提供HTTP客户端和服务器的实现。

Get, Head, Post, 和 PostForm 发出 HTTP (or HTTPS) 请求:

几个例子如下：

```
resp, err := http.Get("http://cpython.org/")
...
resp, err := http.Post("http://cpython.org/upload", "image/jpeg", &buf)
...
resp, err := http.PostForm("http://cpython.org/form",
	url.Values{"key": {"Value"}, "id": {"123"}})
```

cpython.org 并没有Post接口，请勿测试！！！


如果client接收完数据，必须关闭client。

```
resp, err := http.Get("http://cpython.org/")
if err != nil {
	// handle error
}
defer resp.Body.Close()  //完成后关闭
body, err := ioutil.ReadAll(resp.Body)
// ...
fmt.Println(body)
```


## get 请求的例子

```
client := &http.Client{
	CheckRedirect: redirectPolicyFunc,
}
// 定义一个client，再请求
resp, err := client.Get("http://example.com")
// ...

// New 一个请求
req, err := http.NewRequest("GET", "http://example.com", nil)
// ... 增加header
req.Header.Add("If-None-Match", `W/"wyzzy"`)
resp, err := client.Do(req) //开始访问
```


TLS配置

```
tr := &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}
client := &http.Client{Transport: tr}
resp, err := client.Get("https://example.com")
```


## server

```
http.Handle("/foo", fooHandler)

http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
})

log.Fatal(http.ListenAndServe(":8080", nil))
```

就看代码，监听了8080端口，等待http连接，根据path绑定了分配了处理函数。


```
s := &http.Server{
	Addr:           ":8080",
	Handler:        myHandler,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())
// 配置了一些参数
```

## 代码分析

http 公共methods

```
const (
    MethodGet     = "GET"
    MethodHead    = "HEAD"
    MethodPost    = "POST"
    MethodPut     = "PUT"
    MethodPatch   = "PATCH" // RFC 5789
    MethodDelete  = "DELETE"
    MethodConnect = "CONNECT"
    MethodOptions = "OPTIONS"
    MethodTrace   = "TRACE"
)
```

    var DefaultServeMux = &defaultServeMux  //路由

```
Handle //
HandleFunc //
注册DefaultServeMux 路由函数，处理path路由。
```


## server 的两个例子

```
//1.
package main

import (
	"io"
	"net/http"
	"log"
)

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	log.Fatal(http.ListenAndServe(":12345", nil))
}
```


```
//2.
import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	log.Fatal(err)
}
```

## Serve

serve 接收Listener上HTTP接入（连接），为每一个连接着建立一个新的接收线程。每一个服务goroutines（线程）读requests并且调用处理函数。如果处理函数是nil，就是用默认的路由（DefaultServeMux）。


func ServeTLS(l net.Listener, handler Handler, certFile, keyFile string) error

ServeTLS 接收HTTPS连接。但是参数需要提供CA证书和私钥（域名的私钥是要去专门的机构申请的）。

## SetCookie
     func SetCookie(w ResponseWriter, cookie *Cookie)

这是服务器端函数，给Response添加一个cookies，cookies的key必须正确，不正确cookie的key会被丢弃。


## Client

```
type Client struct {
    // Transport specifies the mechanism by which individual
    // HTTP requests are made.
    // If nil, DefaultTransport is used.
    Transport RoundTripper

    //这个部分，我查了一些资料，看了文档，实话讲，我并没有吃透这个部分。
    //我只记得有代理部分。其他未GET到其中的含义。希望以后我有机会来补充这个。2018.08.11

    // CheckRedirect specifies the policy for handling redirects.
    // If CheckRedirect is not nil, the client calls it before
    // following an HTTP redirect. The arguments req and via are
    // the upcoming request and the requests made already, oldest
    // first. If CheckRedirect returns an error, the Client's Get
    // method returns both the previous Response (with its Body
    // closed) and CheckRedirect's error (wrapped in a url.Error)
    // instead of issuing the Request req.
    // As a special case, if CheckRedirect returns ErrUseLastResponse,
    // then the most recent response is returned with its body
    // unclosed, along with a nil error.
    //
    // If CheckRedirect is nil, the Client uses its default policy,
    // which is to stop after 10 consecutive requests.
    CheckRedirect func(req *Request, via []*Request) error

    // 如果请求服务器返回是30x，会使用这个参数来解决redirect请求问题
    // Jar specifies the cookie jar.
    //
    // The Jar is used to insert relevant cookies into every
    // outbound Request and is updated with the cookie values
    // of every inbound Response. The Jar is consulted for every
    // redirect that the Client follows.
    //
    // If Jar is nil, cookies are only sent if they are explicitly
    // set on the Request.
    Jar CookieJar

    // 可以使用 setcookies来这是 cookie
    // Cookie操作server和client都有涉及到。

    // Timeout specifies a time limit for requests made by this
    // Client. The timeout includes connection time, any
    // redirects, and reading the response body. The timer remains
    // running after Get, Head, Post, or Do return and will
    // interrupt reading of the Response.Body.
    //
    // A Timeout of zero means no timeout.
    //
    // The Client cancels requests to the underlying Transport
    // using the Request.Cancel mechanism. Requests passed
    // to Client.Do may still set Request.Cancel; both will
    // cancel the request.
    //
    // For compatibility, the Client will also use the deprecated
    // CancelRequest method on Transport if found. New
    // RoundTripper implementations should use Request.Cancel
    // instead of implementing CancelRequest.
    Timeout time.Duration


}

```

Client是一个http 客户端，DefaultClient 是一个使用DefaultTransport可使用客户端。
Client的Transport 是具有内部标示的（缓存了TCP的连接），所以Clients是可以被重用（reused）的，而不是采取多次新建。Client可以安全地同时使用多个线程（goroutine）。

Client在RoundTripper（Transport）上增加了http的一些细节处理，例如：Cookies和redirects。

## Do
     func (c *Client) Do(req *Request) (*Response, error)

Do 发送一个HTTP 请求并且返回HTTP response（响应）。 并按照客户端上配置的策略（例如redirects， cookies，auth）返回HTTP response。

如果由客户端策略（如CheckRedirect）引起，或者无法说出HTTP（如网络连接问题），则返回错误。错误不是指：非2xx状态码。（状态码404，也是一种正确的返回值）

如果返回的错误为nil，则响应将包含用户应关闭的非零主体。如果Body未关闭，则客户端的基础RoundTripper（通常为Transport）可能无法重新使用到服务器的持久TCP连接以用于后续“保持活动”请求。

请求Body（如果非零）将由底层传输关闭，即使出现错误也是如此。

出错时，可以忽略任何响应。只有在CheckRedirect失败时才会出现带有非零错误的非零响应，即使这样，返回的Response.Body也已关闭。

通常使用Get，Post或PostForm代替Do.

## GET
    func (c *Client) Get(url string) (resp *Response, err error)

获取GET到指定URL的内容。如果响应是以下重定向代码之一，则在调用客户端的CheckRedirect函数后，Get将跟随重定向：

      301（永久移动）
      302（找到）
      303（见其他）
      307（临时重定向）
      308（永久重定向）

如果客户端的CheckRedirect功能失败或者存在HTTP协议错误，则会返回错误。非2xx响应不会导致错误。

当err为nil时，resp总是包含一个非零的resp.Body。调用者在完成阅读后应该关闭resp.Body。

要使用自定义标头发出请求，请使用NewRequest和Client.Do。


## Head

    func (c *Client) Head(url string) (resp *Response, err error)

Head向指定的URL发出HEAD。

## Post
    func (c *Client) Post(url string, contentType string, body io.Reader) (resp *Response, err error)

发布POST到指定的URL。

调用者在完成read后应该关闭resp.Body。

如果提供的主体是io.Closer，则在请求后将其关闭。

要设置自定义标头，请使用NewRequest和Client.Do。

有关如何处理重定向的详细信息，请参阅Client.Do方法文档。

## PostForm

    func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)

PostForm向指定的URL发出POST，数据的键和值URL编码为请求主体。

Content-Type标头设置为application / x-www-form-urlencoded。要设置其他标头，请使用NewRequest和Client.Do。

当err为nil时，resp总是包含一个非零的resp.Body。调用者在完成阅读后应该关闭resp.Body。

有关如何处理重定向的详细信息，请参阅Client.Do方法文档。

## Cookie

```
type Cookie struct {
    Name  string
    Value string

    Path       string    // optional
    Domain     string    // optional
    Expires    time.Time // optional
    RawExpires string    // for reading cookies only

    // MaxAge=0 means no 'Max-Age' attribute specified.
    // MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
    // MaxAge>0 means Max-Age attribute present and given in seconds
    MaxAge   int
    Secure   bool
    HttpOnly bool
    Raw      string
    Unparsed []string // Raw text of unparsed attribute-value pairs
}
```

Cookie表示在HTTP响应的Set-Cookie标头或HTTP请求的Cookie标头中发送的HTTP Cookie。


## CookieJar

```
type CookieJar interface {
    // SetCookies handles the receipt of the cookies in a reply for the
    // given URL.  It may or may not choose to save the cookies, depending
    // on the jar's policy and implementation.
    SetCookies(u *url.URL, cookies []*Cookie)

    // Cookies returns the cookies to send in a request for the given URL.
    // It is up to the implementation to honor the standard cookie use
    // restrictions such as in RFC 6265.
    Cookies(u *url.URL) []*Cookie
}
```

CookieJar管理HTTP请求中cookie的存储和使用。

CookieJar的实现必须安全，以便多个goroutine并发使用。

net / http / cookiejar包提供了一个CookieJar实现。


## Request

```
type Request struct {
    //
    // HTTP方法（GET，POST，PUT等）。 默认是GET
    //
    //
    // Go's HTTP 客户端请求不支持CONNECT，
    Method string


    // 对于服务器端URL，就是URI
    // 对于客户端请求的时候，URL就是要请求的URL地址
    URL *url.URL

    // The protocol version for incoming server requests.
    //
    // For client requests these fields are ignored. The HTTP
    // client code always uses either HTTP/1.1 or HTTP/2.
    // See the docs on Transport for details.
    Proto      string // "HTTP/1.0"
    ProtoMajor int    // 1
    ProtoMinor int    // 0

    // Header contains the request header fields either received
    // by the server or to be sent by the client.
    //
    // If a server received a request with header lines,
    //
    //	Host: example.com
    //	accept-encoding: gzip, deflate
    //	Accept-Language: en-us
    //	fOO: Bar
    //	foo: two
    //
    // then
    //
    //	Header = map[string][]string{
    //		"Accept-Encoding": {"gzip, deflate"},
    //		"Accept-Language": {"en-us"},
    //		"Foo": {"Bar", "two"},
    //	}
    //
    // For incoming requests, the Host header is promoted to the
    // Request.Host field and removed from the Header map.
    //
    // HTTP defines that header names are case-insensitive. The
    // request parser implements this by using CanonicalHeaderKey,
    // making the first character and any characters following a
    // hyphen uppercase and the rest lowercase.
    //
    // For client requests, certain headers such as Content-Length
    // and Connection are automatically written when needed and
    // values in Header may be ignored. See the documentation
    // for the Request.Write method.

    // 是一个map[string][]string,key是字符串，v是字符串数组
    Header Header

    // Body 请求的正文.
    //
    // 对于客户端可以没有body

    // 服务器端必须返回body，没有就填写EOF，请求结束。
    Body io.ReadCloser

    // GetBody  客户端调用返回数据
    //
    // server unused.
    GetBody func() (io.ReadCloser, error)

    // ContentLength 内容的长度 .
    //  -1 表示不确定.
    //  >= 0 表示可以从Body中读取这些长度的内容。

    ContentLength int64

    // TransferEncoding lists the transfer encodings from outermost to
    // innermost. An empty list denotes the "identity" encoding.
    // TransferEncoding can usually be ignored; chunked encoding is
    // automatically added and removed as necessary when sending and
    // receiving requests.
    TransferEncoding []string

    // Close indicates whether to close the connection after
    // replying to this request (for servers) or after sending this
    // request and reading its response (for clients).
    //
    // For server requests, the HTTP server handles this automatically
    // and this field is not needed by Handlers.
    //
    // For client requests, setting this field prevents re-use of
    // TCP connections between requests to the same hosts, as if
    // Transport.DisableKeepAlives were set.
    Close bool

    // For server requests Host specifies the host on which the
    // URL is sought. Per RFC 2616, this is either the value of
    // the "Host" header or the host name given in the URL itself.
    // It may be of the form "host:port". For international domain
    // names, Host may be in Punycode or Unicode form. Use
    // golang.org/x/net/idna to convert it to either format if
    // needed.
    //
    // For client requests Host optionally overrides the Host
    // header to send. If empty, the Request.Write method uses
    // the value of URL.Host. Host may contain an international
    // domain name.
    Host string

    // get 请求的时候是 参数
    // POST请求的时候是form 数据
    Form url.Values

    // PostForm contains the parsed form data from POST, PATCH,
    // or PUT body parameters.
    //
    // This field is only available after ParseForm is called.
    // The HTTP client ignores PostForm and uses Body instead.
    // server 中使用
    PostForm url.Values

    // MultipartForm is the parsed multipart form, including file uploads.
    // This field is only available after ParseMultipartForm is called.
    // The HTTP client ignores MultipartForm and uses Body instead.
    // server中使用
    MultipartForm *multipart.Form

    // Trailer specifies additional headers that are sent after the request
    // body.
    //
    // For server requests the Trailer map initially contains only the
    // trailer keys, with nil values. (The client declares which trailers it
    // will later send.)  While the handler is reading from Body, it must
    // not reference Trailer. After reading from Body returns EOF, Trailer
    // can be read again and will contain non-nil values, if they were sent
    // by the client.
    //
    // For client requests Trailer must be initialized to a map containing
    // the trailer keys to later send. The values may be nil or their final
    // values. The ContentLength must be 0 or -1, to send a chunked request.
    // After the HTTP request is sent the map values can be updated while
    // the request body is read. Once the body returns EOF, the caller must
    // not mutate Trailer.
    //
    // Few HTTP clients, servers, or proxies support HTTP trailers.
    Trailer Header

    // RemoteAddr allows HTTP servers and other software to record
    // the network address that sent the request, usually for
    // logging. This field is not filled in by ReadRequest and
    // has no defined format. The HTTP server in this package
    // sets RemoteAddr to an "IP:port" address before invoking a
    // handler.
    // This field is ignored by the HTTP client.
    RemoteAddr string

    // RequestURI is the unmodified Request-URI of the
    // Request-Line (RFC 2616, Section 5.1) as sent by the client
    // to a server. Usually the URL field should be used instead.
    // It is an error to set this field in an HTTP client request.
    RequestURI string

    // TLS allows HTTP servers and other software to record
    // information about the TLS connection on which the request
    // was received. This field is not filled in by ReadRequest.
    // The HTTP server in this package sets the field for
    // TLS-enabled connections before invoking a handler;
    // otherwise it leaves the field nil.
    // This field is ignored by the HTTP client.
    TLS *tls.ConnectionState

    // Cancel is an optional channel whose closure indicates that the client
    // request should be regarded as canceled. Not all implementations of
    // RoundTripper may support Cancel.
    //
    // For server requests, this field is not applicable.
    //
    // Deprecated: Use the Context and WithContext methods
    // instead. If a Request's Cancel field and context are both
    // set, it is undefined whether Cancel is respected.
    Cancel <-chan struct{}

    // Response is the redirect response which caused this request
    // to be created. This field is only populated during client
    // redirects.
    Response *Response
    // contains filtered or unexported fields
}

```


##  NewRequest

   func NewRequest(method, url string, body io.Reader) (*Request, error)

NewRequest在给定方法，URL和可选主体的情况下返回新请求。

如果提供的正文也是io.Closer，则返回的Request.Body设置为body，将由客户端方法Do，Post和PostForm以及Transport.RoundTrip关闭。

NewRequest返回适用于Client.Do或Transport.RoundTrip的请求。要创建用于测试服务器处理程序的请求，请使用net / http / httptest包中的NewRequest函数，使用ReadRequest，或手动更新Request字段。有关入站和出站请求字段之间的区别，请参阅请求类型的文档。

如果body的类型为* bytes.Buffer，* bytes.Reader或* strings.Reader，则返回的请求的ContentLength设置为其精确值（而不是-1），GetBody将被填充（因此307和308重定向可以重放如果ContentLength为0，则Body设置为NoBody。


## AddCookie
    func（r * Request）AddCookie（c * Cookie）

AddCookie为请求添加了一个cookie。根据RFC 6265第5.4节，AddCookie不会附加多个Cookie标头字段。这意味着所有cookie（如果有的话）都写入同一行，用分号分隔。

## Cookies， Cookie

    func (r *Request) Cookies() []*Cookie
    func (r *Request) Cookie(name string) (*Cookie, error)



## form 表单或者文件传输

     func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)

     func (r *Request) FormValue(key string) string


##  Response


```
 type Response struct {
    Status     string // e.g. "200 OK"
    StatusCode int    // e.g. 200
    Proto      string // e.g. "HTTP/1.0"
    ProtoMajor int    // e.g. 1
    ProtoMinor int    // e.g. 0

    // Header maps header keys to values. If the response had multiple
    // headers with the same key, they may be concatenated, with comma
    // delimiters.  (Section 4.2 of RFC 2616 requires that multiple headers
    // be semantically equivalent to a comma-delimited sequence.) When
    // Header values are duplicated by other fields in this struct (e.g.,
    // ContentLength, TransferEncoding, Trailer), the field values are
    // authoritative.
    //
    // Keys in the map are canonicalized (see CanonicalHeaderKey).
    Header Header

    // Body represents the response body.
    //
    // The response body is streamed on demand as the Body field
    // is read. If the network connection fails or the server
    // terminates the response, Body.Read calls return an error.
    //
    // The http Client and Transport guarantee that Body is always
    // non-nil, even on responses without a body or responses with
    // a zero-length body. It is the caller's responsibility to
    // close Body. The default HTTP client's Transport may not
    // reuse HTTP/1.x "keep-alive" TCP connections if the Body is
    // not read to completion and closed.
    //
    // The Body is automatically dechunked if the server replied
    // with a "chunked" Transfer-Encoding.
    Body io.ReadCloser

    // ContentLength records the length of the associated content. The
    // value -1 indicates that the length is unknown. Unless Request.Method
    // is "HEAD", values >= 0 indicate that the given number of bytes may
    // be read from Body.
    ContentLength int64

    // Contains transfer encodings from outer-most to inner-most. Value is
    // nil, means that "identity" encoding is used.
    TransferEncoding []string

    // Close records whether the header directed that the connection be
    // closed after reading Body. The value is advice for clients: neither
    // ReadResponse nor Response.Write ever closes a connection.
    Close bool

    // Uncompressed reports whether the response was sent compressed but
    // was decompressed by the http package. When true, reading from
    // Body yields the uncompressed content instead of the compressed
    // content actually set from the server, ContentLength is set to -1,
    // and the "Content-Length" and "Content-Encoding" fields are deleted
    // from the responseHeader. To get the original response from
    // the server, set Transport.DisableCompression to true.
    Uncompressed bool
    // 是否压缩

    // Trailer maps trailer keys to values in the same
    // format as Header.
    //
    // The Trailer initially contains only nil values, one for
    // each key specified in the server's "Trailer" header
    // value. Those values are not added to Header.
    //
    // Trailer must not be accessed concurrently with Read calls
    // on the Body.
    //
    // After Body.Read has returned io.EOF, Trailer will contain
    // any trailer values sent by the server.
    Trailer Header

    // Request is the request that was sent to obtain this Response.
    // Request's Body is nil (having already been consumed).
    // This is only populated for Client requests.
    Request *Request

    // TLS contains information about the TLS connection on which the
    // response was received. It is nil for unencrypted responses.
    // The pointer is shared between responses and should not be
    // modified.
    TLS *tls.ConnectionState
}
```


##  Transfer 例子

```
var DefaultTransport RoundTripper = &Transport{
    Proxy: ProxyFromEnvironment,
    DialContext: (&net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
        DualStack: true,
    }).DialContext,
    MaxIdleConns:          100,
    IdleConnTimeout:       90 * time.Second,
    TLSHandshakeTimeout:   10 * time.Second,
    ExpectContinueTimeout: 1 * time.Second,
}
```

原文档后面有sever部分。此次，我没有写。因为我要写几个httpclient的例子了。
