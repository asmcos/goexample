package main

import (
	"io/ioutil"
	"net/http"
	"fmt"
)

func main (){

	resp, err := http.Get("https://cn.bing.com/")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// ...
	_ = string(body)
	fmt.Println(resp)

	return 
}
