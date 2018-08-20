package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {

	//creating the proxyURL
	/*
		proxyStr := "http://localhost:8888"
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			log.Println(err)
		}
	*/
	//creating the URL to be loaded through the proxy
	urlStr := "http://cn.bing.com"
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
	}

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		// Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	//generating the HTTP GET request
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Println(err)
	}
	request.Header.Set("User-Agent", "Golang requests")
	//request.Header.Set("Accept-Encoding", "gzip")
	//calling the URL
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}

	//getting the response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	//printing the response
	log.Println(string(data))
	fmt.Println(response)
}
