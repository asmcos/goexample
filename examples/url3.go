package main

import (
	"fmt"
	"net/url"
)

func main(){

	u := new(url.URL)

	u.Scheme = "https"
	u.Host = "google.com"
	q := u.Query()
	q.Set("q", "golang")
	u.RawQuery = q.Encode()
	fmt.Println(u)

}
