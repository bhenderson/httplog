package httplog

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func startServer(t *testing.T) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	return httptest.NewServer(h)
}
