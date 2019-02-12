package httplog

import (
	"io"
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
}

func NewLogger(r *http.Request) *Logger {
	l := &Logger{}
	l.Request.Method = r.Method
	l.Request.URL = r.URL
	l.Request.Header = r.Header
	l.Start = time.Now()
	r.Body, l.Request.Body = copyBody(r.Body)

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
		if f, ok := resp.Request.Context().Value(ContextFormat).(string); ok {
			format = f
		}
	}

	if w == nil {
		w = os.Stdout
	}
	Template.ExecuteTemplate(w, format, l)
}

func (l *Logger) Duration() time.Duration {
	return l.End.Sub(l.Start)
}
