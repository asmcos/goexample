package requests

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	//"fmt"
)

var VERSION string = "0.2"

type request struct {
	httpreq *http.Request
	Header  *http.Header
	Client  *http.Client
}

type response struct {
	httpresp *http.Response
	content  []byte
	text     string
}

type Header map[string]string
type Params map[string]string

// {username,password}
type Auth []string

func Requests() *request {

	req := new(request)

	req.httpreq = &http.Request{
		Method:     "GET",
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req.Header = &req.httpreq.Header
	req.httpreq.Header.Set("User-Agent", "Go-Requests "+VERSION)

	req.Client = &http.Client{}

	return req
}

func Get(origurl string, args ...interface{}) (resp *response) {
	req := Requests()

	// call request Get
	resp = req.Get(origurl, args...)
	return resp
}

func (req *request) Get(origurl string, args ...interface{}) (resp *response) {
	// set params ?a=b&b=c
	//set Header
	params := []map[string]string{}

	for _, arg := range args {
		switch a := arg.(type) {
		// arg is Header , set to request header
		case Header:

			for k, v := range a {
				req.Header.Set(k, v)
			}
			// arg is "GET" params
			// ?title=website&id=1860&from=login
		case Params:
			params = append(params, a)
		case Auth:
			// a{username,password}
			req.httpreq.SetBasicAuth(a[0],a[1])
		}
	}

	disturl, _ := buildURLParams(origurl, params...)

	//prepare to Do
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil
	}
	req.httpreq.URL = URL

	res, err := req.Client.Do(req.httpreq)

	if err != nil {
		return nil
	}

	resp = &response{}
	resp.httpresp = res
	return resp
}

// handle URL params
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

/**************/
func (resp *response) Content() []byte {

	defer resp.httpresp.Body.Close()
	var err error
	resp.content, err = ioutil.ReadAll(resp.httpresp.Body)
	if err != nil {
		return nil
	}

	return resp.content
}

func (resp *response) Text() string {
	if resp.content == nil {
		resp.Content()
	}
	resp.text = string(resp.content)
	return resp.text
}
