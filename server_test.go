package httplog

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var output = new(bytes.Buffer)

func startServer(t *testing.T) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	return httptest.NewServer(h)
}
