package httplog

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Request is mostly because http.Request is not JSON marshalable
type Request struct {
	Method string
	URL    *url.URL
	Header http.Header
	Body   string
}

// Response is mostly because http.Response is not JSON marshalable
type Response struct {
	StatusCode int
	Header     http.Header
	Body       string
}

// Logger is the struct passed to the template for logging
type Logger struct {
	Request  Request
	Response Response

	Start time.Time
	End   time.Time
	Error error

	req  *http.Request
	resp *http.Response
}

func NewLogger(r *http.Request) *Logger {
	l := &Logger{}
	l.Request.Method = r.Method
	l.Request.URL = r.URL
	l.Request.Header = r.Header
	l.Start = time.Now()
	r.Body, l.Request.Body = copyBody(r.Body)
	l.req = r

	RequestCtx(r, ContextLogger, l)

	return l
}

func (l *Logger) Log(w io.Writer, resp *http.Response, err error) {
	format := DefaultFormat

	l.End = time.Now()
	if err != nil {
		l.Error = err
	}
	if resp != nil {
		l.Response.StatusCode = resp.StatusCode
		l.Response.Header = resp.Header
		resp.Body, l.Response.Body = copyBody(resp.Body)
		l.resp = resp
		ctx := resp.Request.Context()
		val := ctx.Value(ContextFormat)
		if f, ok := val.(string); ok {
			format = f
		}
	}

	if w == nil {
		w = os.Stdout
	}
	err = Template.ExecuteTemplate(w, format, l)
	if err != nil {
		log.Println(err)
	}
}

func (l *Logger) Duration() time.Duration {
	return l.End.Sub(l.Start)
}

func (l *Logger) RawRequest() *http.Request {
	return l.req
}

func (l *Logger) RawResponse() *http.Response {
	return l.resp
}
