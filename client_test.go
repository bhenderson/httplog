package httplog

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	DefaultTransport.Output = output
	http.DefaultTransport = DefaultTransport

	s := startServer(t)
	defer s.Close()

	// KV test
	t.Run("KV", func(t *testing.T) {
		output.Reset()
		req, _ := http.NewRequest("GET", s.URL, nil)
		RequestCtx(req, ContextFormat, "KV")
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, dumpBody(resp.Body))
		l := resp.Request.Context().Value(ContextLogger).(*Logger)
		exp := "method=GET url=" + s.URL + " code=200 duration=" + l.Duration().String() + "\n"
		act := output.String()
		assert.Equal(t, exp, act)
	})

	t.Run("JSON", func(t *testing.T) {
		output.Reset()
		req, _ := http.NewRequest("GET", s.URL, nil)
		RequestCtx(req, ContextFormat, "JSON")
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, dumpBody(resp.Body))
		l := resp.Request.Context().Value(ContextLogger).(*Logger)
		exp, _ := json.Marshal(l)
		act := output.String()
		assert.Equal(t, string(exp)+"\n", act)
	})

	t.Run("Body", func(t *testing.T) {
		content := `{"hello": "world"}`
		body := bytes.NewBufferString(content)
		req, _ := http.NewRequest("POST", s.URL, body)
		RequestCtx(req, ContextFormat, "JSON")
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, dumpBody(resp.Body))
		l := resp.Request.Context().Value(ContextLogger).(*Logger)

		t.Run("Request", func(t *testing.T) {
			exp := content
			act := l.Request.Body
			assert.Equal(t, exp, act)
		})

		t.Run("Response", func(t *testing.T) {
			exp := "hello world"
			act := l.Response.Body
			assert.Equal(t, exp, act)
		})
	})
}

func dumpBody(r io.ReadCloser) string {
	defer r.Close()
	content, _ := ioutil.ReadAll(r)
	return string(content)
}
