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
	fmt.Println(u.RawPath)   //原始存储
	fmt.Println(u.String())

}
