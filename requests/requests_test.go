package requests

import (
	"testing"
)

func TestNew(t *testing.T) {
	req := Requests()

	Requests().Get("http://www.cpython.org", Params{"a": "b"})

	req.Get("http://go.xiulian.net.cn", Header{"Content-length": "0"}, Params{"c": "d", "e": "f"}, Params{"c": "a"})

	Get("http://jeapedu.com")

}
