package requests

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type request struct {
	httpreq *http.Request
}

type response struct {
	httpresp *http.Response
}

type Header map[string]string
type Params map[string]string

func Requests() *request {

	req := new(request)

	req.httpreq = &http.Request{
		Method:     "GET",
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	req.httpreq.Header.Set("User-Agent", "Go-Requests")

	return req
}

func Get(origurl string, args ...interface{}) {
	req := Requests()

	req.Get(origurl, args)
}

func (req *request) Get(origurl string, args ...interface{}) {
	// set params ?a=b&b=c
	//set Header
	params := []map[string]string{}
	for _, arg := range args {
		switch a := arg.(type) {
		case Header:
			fmt.Println("header")
			fmt.Println(a)
		case Params:
			params = append(params, a)
		}
	}

	//
	disturl, _ := buildURLParams(origurl, params...)

	fmt.Println(disturl)
}

func buildURLParams(userURL string, params ...map[string]string) (string, error) {
	parsedURL, err := url.Parse(userURL)

	if err != nil {
		return "", err
	}

	parsedQuery, err := url.ParseQuery(parsedURL.RawQuery)

	if err != nil {
		return "", nil
	}

	for _, param := range params {
		for key, value := range param {
			parsedQuery.Add(key, value)
		}
	}
	return addQueryParams(parsedURL, parsedQuery), nil
}

func addQueryParams(parsedURL *url.URL, parsedQuery url.Values) string {
	if len(parsedQuery) > 0 {
		return strings.Join([]string{strings.Replace(parsedURL.String(), "?"+parsedURL.RawQuery, "", -1), parsedQuery.Encode()}, "?")
	}
	return strings.Replace(parsedURL.String(), "?"+parsedURL.RawQuery, "", -1)
}
