package requests

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	// example 1
	req := Requests()

	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Get("http://go.xiulian.net.cn", Header{"Content-Length": "0"}, Params{"c": "d", "e": "f"}, Params{"c": "a"})

	// example 2
	h := Header{
		"Referer":         "http://www.jeapedu.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
	}

	Get("http://jeapedu.com", h, Header{"Content-Length": "1024"})

	// example 3
	p := Params{
		"title": "The blog",
		"name":  "file",
		"id":    "12345",
	}
	resp := Requests().Get("http://www.cpython.org", p)

	resp.Text()
	fmt.Println(resp.text)
}
