package curl_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bhenderson/httplog"
	"github.com/bhenderson/httplog/curl"
)

func Example() {
	httplog.Trace(curl.Format)

	body := bytes.NewBufferString("hello")
	req, _ := http.NewRequest("POST", "http://example.com", body)
	l := httplog.NewLogger(req)

	// would be http.DefaultClient.Do(req), but it's hard to write an example using http
	resp := newResponse()
	l.Log(os.Stdout, resp, nil)
	// Output: curl -X 'POST' -d 'hello' 'http://example.com'
	// HTTP/1.0 200 OK

	// Host: example.com

	// Content-Type: text/plain

	// world
}

func newResponse() *http.Response {
	return &http.Response{
		StatusCode: 200,
		Status:     "OK",
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header: http.Header{
			"Content-Type": []string{"text/plain"},
			"Host":         []string{"example.com"},
		},
		Body:    ioutil.NopCloser(bytes.NewBufferString("world")),
		Request: req,
	}

}
