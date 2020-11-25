package worker 

import (
	"bytes"
	"net/http"
	"io/ioutil"

	"github.com/json-iterator/go"
)


type Response struct {
	Status			string
	StatusCode 		int
	Proto			string
	Header 			http.Header
	Body 			[]byte	
}

type Worker struct {
	Client 			*http.Client
	Header 			http.Header
	Cookies 		map[string][]*http.Cookie
}

func (w Worker) Get(gurl string, setCookies bool) (*Response, error) {

	// make req
	req, _ := http.NewRequest("GET", gurl, nil)
	if w.Header != nil {
		req.Header = w.Header
	}
	if cookie, ok := w.Cookies[gurl]; ok{
		for _,c := range cookie {
			req.AddCookie(c)
		}
	}

	// do request
	resp, err := w.Client.Do(req)
	if err != nil {
		return  nil, err
	}

	// read body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// cookie
	if setCookies {
		cookies := resp.Cookies()
		w.SetCookies(gurl, cookies)
	}

	// set resp
	given := &Response{
		Status 		: resp.Status,
		StatusCode 	: resp.StatusCode,
		Proto		: resp.Proto,
		Header 		: resp.Header,
		Body 		: body,
	}

	return given, nil
}

func (w Worker) Post(gurl string, setCookies bool, data interface{}) (*Response, error) {
	// make req
	b, err := jsoniter.Marshal(data)
	if err != nil{
		return nil, err
	}
	req, _ := http.NewRequest( "POST", gurl, bytes.NewBuffer(b))
	if w.Header != nil {
		req.Header = w.Header
	}
	if cookie, ok := w.Cookies[gurl]; ok{
		for _,c := range cookie {
			req.AddCookie(c)
		}
	}

	// do req
	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// cookie
	if setCookies{
		cookies := resp.Cookies()
		w.SetCookies(gurl, cookies)	
	}

	given := &Response{
		Status 		: resp.Status,
		StatusCode 	: resp.StatusCode,
		Proto		: resp.Proto,
		Header 		: resp.Header,
		Body 		: body,
	}

	return given, nil
}

func (w Worker) SetCookies (url string, cookies []*http.Cookie){
	if given, ok := w.Cookies[url]; ok {
		for _,c := range cookies {
			w.Cookies[url] = append( given, c)
		}
	} else {
		w.Cookies[url] = cookies
	}
}

func InitWorker() *Worker{
	c := &http.Client{}
	w := &Worker{
		Client:c,
		Cookies: make(map[string][]*http.Cookie),
	}
	return w
}



















