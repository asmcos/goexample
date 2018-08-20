package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Experiment struct {
	URL          string
	ExtraHeaders map[string]string
}

func (e Experiment) Describe() string {
	s := e.URL
	for k, v := range e.ExtraHeaders {
		s += fmt.Sprintf(" %s=%s", k, v)
	}
	return s
}

func main() {
	/*	addAcceptEncodingHeader := map[string]string{
		"Accept-Encoding": "gzip",
	}*/

	experiments := []Experiment{
		{"http://httpbin.org/", nil},
	}

	for _, e := range experiments {
		fmt.Printf("==== %s ====\n", e.Describe())

		req, err := http.NewRequest(http.MethodGet, e.URL, nil)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range e.ExtraHeaders {
			req.Header.Add(k, v)
		}

		err = DebugRequest(req)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("")
	}
}

func DebugRequest(req *http.Request) error {
	b, err := httputil.DumpRequestOut(req, false)
	if err != nil {
		return err
	}

	proxyStr := "http://localhost:8888"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(b))

	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//adding the Transport object to the http Client
	client := &http.Client{
		Transport: transport,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	b, err = httputil.DumpResponse(res, false)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	fmt.Println(res)

	fmt.Println("Uncompressed", res.Uncompressed)

	return nil
}
