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
