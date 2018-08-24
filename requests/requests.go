package requests

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"fmt"
	"compress/gzip"
	"encoding/json"
	"os"
	"crypto/tls"

)

var VERSION string = "0.3"

type request struct {
	httpreq *http.Request
	Header  *http.Header
	Client  *http.Client
	Debug   int
}

type response struct {
	httpresp *http.Response
	content  []byte
	text     string
	req * request
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


	req.RequestDebug()


	res, err := req.Client.Do(req.httpreq)

	if err != nil {
		fmt.Println(err)
		return nil
	}



	resp = &response{}
	resp.httpresp = res
	resp.req = req

	resp.ResponseDebug()
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

func (req *request) RequestDebug(){


  if req.Debug != 1{
		return
	}

	fmt.Println("===========Go RequestDebug libaray ============")

	message, err := httputil.DumpRequestOut(req.httpreq, false)
	if err != nil {
		return
	}
  fmt.Println(string(message))

}

func (req * request) Proxy(proxyurl string){

	urli := url.URL{}
	urlproxy, err:= urli.Parse(proxyurl)
  if err != nil {
		fmt.Println("Set proxy failed")
		return
	}
	req.Client.Transport = &http.Transport{
		Proxy:http.ProxyURL(urlproxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

}

/**************/
func (resp *response) ResponseDebug(){


	if resp.req.Debug != 1 {
		return
	}

	fmt.Println("===========Go ResponseDebug libaray ============")

	message, err := httputil.DumpResponse(resp.httpresp, false)
	if err != nil {
		return
	}

	fmt.Println(string(message))
}

func (resp *response) Content() []byte {

	defer resp.httpresp.Body.Close()
	var err error

  var Body = resp.httpresp.Body
	if resp.httpresp.Header.Get("Content-Encoding") == "gzip" && resp.req.Header.Get("Accept-Encoding") != "" {
		// fmt.Println("gzip")
		reader, err := gzip.NewReader(Body)
		if err != nil {
			return nil
		}
		Body = reader
	}

	resp.content, err = ioutil.ReadAll(Body)
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

func (resp *response) SaveFile(filename string) error {
	if resp.content == nil {
		resp.Content()
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(resp.content)
	f.Sync()

	return err
}

func (resp *response) Json(v interface{}) error {
	if resp.content == nil {
		resp.Content()
	}
	return json.Unmarshal(resp.content, v)
}

/*******/
// I don't want import fmt or delete fmt replay
func empty (){
	_ = fmt.Println
}
