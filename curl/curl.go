package curl

import (
	"text/template"

	"github.com/bhenderson/httplog"
	"github.com/moul/http2curl"
)

// Format is the template name
var Format = "CURL"

// Template adds to httplog.Template and is named only for documentation
var Template = template.Must(httplog.Template.
	Funcs(template.FuncMap{
		"curl": http2curl.GetCurlCommand,
	}).
	Parse(
		`
{{- define "CURL" -}}
{{ curl .RawRequest }}
{{ dumpout .RawResponse true | printf "%s" }}
{{ end -}}
`,
	),
)
