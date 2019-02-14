package curl

import (
	"regexp"
	"sort"
	"testing"

	"github.com/bhenderson/httplog"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	reg := regexp.MustCompile(`; defined templates are: "(.+)", "(.+)", "(.+)", "(.+)"`)
	names := httplog.Template.DefinedTemplates()

	{
		assert.Regexp(t, reg, names)
	}

	{
		exp := []string{
			"CURL",
			"JSON",
			"KV",
			"httplog",
		}
		act := reg.FindStringSubmatch(names)
		sort.Strings(act)
		assert.Equal(t, exp, act[1:])
	}
}
