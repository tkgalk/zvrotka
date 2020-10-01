package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Customer HTTP client with 30 seconds timeout.
var netClient = &http.Client{
	Timeout: time.Second * 30,
}

// Request is a wrapper object generated from the requests.http file.
type Request struct {
	Name    string
	Method  string
	Url     string
	Headers []string
	Body    string
}

// Response is a wrapper object for the http.Response.
type Response struct {
	Res      *http.Response
	Duration time.Duration
}

func (r Request) String() string {
	return fmt.Sprintf("{\nName: %s\nMethod: %s\nURL: %s\nHeaders: %v\nBody: %s\n}", r.Name, r.Method, r.Url, r.Headers, r.Body)
}

// SendRequest makes an HTTP call and returns a Response object.
func SendRequest(req Request) (*Response, error) {
	parsedUrl, err := url.Parse(req.Url)
	if err != nil {
		return nil, err
	}

	body := ioutil.NopCloser(strings.NewReader(req.Body))
	r := &http.Request{
		Method: req.Method,
		URL:    parsedUrl,
		Header: map[string][]string{},
		Body:   body,
	}

	start := time.Now()
	res, err := netClient.Do(r)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start)

	result := &Response{
		Res:      res,
		Duration: duration,
	}

	return result, nil
}
