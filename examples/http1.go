package main

import (
	"io/ioutil"
	"net/http"
	"fmt"
)

func main (){

	resp, err := http.Get("http://cpython.org/")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// ...
	fmt.Println(string(body[:]))

	return 
}
