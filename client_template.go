package httplog

import (
	"encoding/json"
	"net/http/httputil"
	"text/template"
	"time"

	"github.com/moul/http2curl"
)

// DefaultFormat is the default template name used for logging. This can be
// changed to any supported or custom template name. The template name can also
// be set in the request context for individual requests.
var DefaultFormat = "KV"

// Template holds the templates that can be used to log request and response
// Note that other funcmaps or custom templates can be added.
var Template = template.Must(template.New("httplog").
	Funcs(template.FuncMap{
		"since":   time.Since,
		"json":    json.Marshal,
		"curl":    http2curl.GetCurlCommand,
		"dumpin":  httputil.DumpRequestOut,
		"dumpout": httputil.DumpResponse,
	}).
	Parse(
		`
{{- define "KV" -}}
{{- if .Error -}}
method={{ .Request.Method }} url={{ printf "%v" .Request.URL }} error={{ printf "%q" .Error }}
{{- else -}}
method={{ .Request.Method }} url={{ printf "%v" .Request.URL }} code={{ .Response.StatusCode }} duration={{ printf "%v" .Duration }}
{{- end }}
{{ end -}}

{{- define "JSON" -}}
{{ json . | printf "%s" }}
{{ end -}}

{{- define "CURL" -}}
{{ curl .Request }}
{{ dumpout .Response true | printf "%s" }}
{{ end -}}
`,
	),
)
