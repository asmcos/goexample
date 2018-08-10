package main

import (

	"fmt"
	"net/url"
)

func main (){
	u, err := url.Parse("http://example.com/path with spaces")
	if err != nil {
    		// log.Fatal(err)
	}
	fmt.Println(u.EscapedPath())	

	fmt.Println(u.Path)
	fmt.Println(u.RawPath)
	fmt.Println(u.RawQuery)
}
