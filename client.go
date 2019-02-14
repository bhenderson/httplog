package httplog

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
)

// DefaultTransport is the default Transport useful for http.DefaultTransport
var DefaultTransport = &Transport{RoundTripper: http.DefaultTransport}

type contextKey struct{ string }

// Context keys
var (
	ContextLogger = &contextKey{"httplog.logger"}
	ContextFormat = &contextKey{"httplog.format"}
)

// Trace sets up http.DefaultTransport to use DefaultTransport with Format named format.
func Trace(format string) {
	if format != "" {
		DefaultFormat = format
	}
	http.DefaultTransport = DefaultTransport
}

// Transport implements http.RoundTripper with logging
type Transport struct {
	http.RoundTripper

	// Set output to the writer you want to log to. Defaults to os.Stdout
	Output io.Writer
}

// RoundTrip implements http.RoundTrip. It also logs the request and response with Logger
func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	l := NewLogger(req)
	defer func() { l.Log(t.Output, resp, err) }()

	return t.RoundTripper.RoundTrip(req)
}

func copyBody(r io.ReadCloser) (io.ReadCloser, string) {
	if r == nil {
		return nil, ""
	}

	defer r.Close()
	content, _ := ioutil.ReadAll(r)
	return ioutil.NopCloser(bytes.NewBuffer(content)), string(content)
}

// RequestCtx is a convenience function for setting the context.
func RequestCtx(r *http.Request, k, v interface{}) {
	req := r.WithContext(context.WithValue(r.Context(), k, v))
	*r = *req
}
