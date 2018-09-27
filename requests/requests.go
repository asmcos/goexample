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
	"net/http/cookiejar"
	"mime/multipart"
	"bytes"
	"io"

)

var VERSION string = "0.4"

type request struct {
	httpreq *http.Request
	Header  *http.Header
	Client  *http.Client
	Debug   int
	Cookies []*http.Cookie
	Forms   url.Values
}

type response struct {
	httpresp *http.Response
	content  []byte
	text     string
	req * request
}

type Header map[string]string
type Params map[string]string
type Datas  map[string]string  // for post form
type Files  map[string]string  // name ,filename

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

  // auto with Cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
			return nil
	}
	req.Client.Jar = jar

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

  //reset Cookies,
	//Client.Do can copy cookie from client.Jar to req.Header
	delete(req.httpreq.Header,"Cookie")

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

	req.ClientSetCookies()

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

	fmt.Println("===========Go RequestDebug ============")

	message, err := httputil.DumpRequestOut(req.httpreq, false)
	if err != nil {
		return
	}
  fmt.Println(string(message))

	if len(req.Client.Jar.Cookies(req.httpreq.URL)) > 0{
		fmt.Println("Cookies:")
		for _, cookie := range req.Client.Jar.Cookies(req.httpreq.URL) {
	            fmt.Println(cookie)
	  }
	}
}

func (req *request ) SetCookie(cookie *http.Cookie){
	req.Cookies = append(req.Cookies,cookie)
}

func (req * request) ClearCookies(){
	req.Cookies = req.Cookies[0:0]
}

func (req * request) ClientSetCookies(){

	if len(req.Cookies) > 0 {
		// 1. Cookies have content, Copy Cookies to Client.jar
		// 2. Clear  Cookies
		req.Client.Jar.SetCookies(req.httpreq.URL, req.Cookies)
		req.ClearCookies()
	}

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

	fmt.Println("===========Go ResponseDebug ============")

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

/**************post*************************/
func Post(origurl string, args ...interface{}) (resp *response) {
	req := Requests()

	// call request Get
	resp = req.Post(origurl, args...)
	return resp
}

func (req *request) Post(origurl string, args ...interface{}) (resp *response) {

	req.httpreq.Method = "POST"
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Forms = url.Values{}

	// set params ?a=b&b=c
	//set Header
	params := []map[string]string{}
	datas := []map[string]string{} // POST
	files := []map[string]string{} //post file

  //reset Cookies,
	//Client.Do can copy cookie from client.Jar to req.Header
	delete(req.httpreq.Header,"Cookie")

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

		case Datas:   //Post form data,packaged in body.
			datas = append(datas,a)
		case Files:
			files = append(files,a)
		case Auth:
			// a{username,password}
			req.httpreq.SetBasicAuth(a[0],a[1])
		}
	}

	disturl, _ := buildURLParams(origurl, params...)

	if len(files) > 0{
		req.buildFilesAndForms(files,datas)

	} else {
			req.buildForms(datas...)
			req.setBodyBytes() // set forms to body
  }
	//prepare to Do
	URL, err := url.Parse(disturl)
	if err != nil {
		return nil
	}
	req.httpreq.URL = URL

	req.ClientSetCookies()

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

func (req * request)setBodyBytes() {

  // maybe
	data := req.Forms.Encode()
	req.httpreq.Body = ioutil.NopCloser(strings.NewReader(data))
	req.httpreq.ContentLength = int64(len(data))
}

func (req * request) buildFilesAndForms(files []map[string]string,datas []map[string]string){

	//handle file multipart

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for _,file := range files{
		for k,v := range file{
			part, err := w.CreateFormFile(k, v)
			if err != nil {
			  fmt.Printf("Upload %s failed!",v)
				panic(err)
			}
			file := openFile(v)
			_, err = io.Copy(part, file)
		}
	}

	for _,data := range datas{
		for k,v := range data{
			w.WriteField(k,v)
		}
	}


  w.Close()
  // set file header example:
	// "Content-Type": "multipart/form-data; boundary=------------------------7d87eceb5520850c",
	req.httpreq.Body = ioutil.NopCloser( bytes.NewReader(b.Bytes()) )
	req.httpreq.ContentLength = int64(b.Len())
	req.Header.Set("Content-Type", w.FormDataContentType())
}

func (req * request)buildForms(datas ...map[string]string)  {

	for _, data := range datas {
		for key, value := range data {
			req.Forms.Add(key, value)
		}
	}

}

func openFile(filename string)*os.File {
    r, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return r
}

/*******/
// I don't want import fmt or delete fmt replay
func empty (){
	_ = fmt.Println
}
