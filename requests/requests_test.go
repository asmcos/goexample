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

  // example 4
  // test authentication usernae,password
	//documentation https://www.httpwatch.com/httpgallery/authentication/#showExample10
	req = Requests()
	resp = req.Get("https://www.httpwatch.com/httpgallery/authentication/authenticatedimage/default.aspx?0.45874470316137206",Auth{"httpwatch","foo"})
	fmt.Println(resp.httpresp)

  
}
