package main

import (
	//"io/ioutil"
	"net/http"
	"fmt"
	"log"
	"crypto/tls"
)

func main (){

//proxyString := "http://127.0.0.1:8888"
//proxyUrl, _ := url.Parse(proxyString)

tr := &http.Transport{
	//Proxy: http.ProxyURL(proxyUrl),
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},     
}

client := &http.Client{}
client.Transport = tr

request, err := http.NewRequest("HEAD", "http://cn.bing.com", nil)
request.Header.Set("User-Agent", "Golang requests")
//request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")

resp, err := client.Do(request)

if err != nil {
    log.Fatalln(err)
    return
}
defer resp.Body.Close()

fmt.Println(resp)


	/*
	resp, err := http.Head("https://cn.bing.com/")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	fmt.Println(resp)
	*/
	return 
}
