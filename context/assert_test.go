package context

import (
	"strconv"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/bilus/scenarigo/assert"
	"github.com/bilus/scenarigo/internal/testutil"
	"github.com/bilus/scenarigo/template"
)

func TestAssertions(t *testing.T) {
	executor := func(r testutil.Reporter, decode func(interface{})) func(testutil.Reporter, interface{}) error {
		var i interface{}
		decode(&i)
		return func(r testutil.Reporter, v interface{}) error {
			a, err := template.Execute(i, map[string]interface{}{
				"assert": assertions,
			})
			if err != nil {
				r.Fatalf("failed to execute template: %s", err)
			}
			return assert.Build(a).Assert(v)
		}
	}
	testutil.RunParameterizedTests(
		t, executor,
		"testdata/assertion/and.yaml",
		"testdata/assertion/or.yaml",
		"testdata/assertion/contains.yaml",
	)
}

func TestLeftArrowFunc(t *testing.T) {
	tests := map[string]struct {
		yaml string
		ok   interface{}
		ng   interface{}
	}{
		"simple": {
			yaml: `'{{f <-}}: 1'`,
			ok:   []int{0, 1},
			ng:   []int{2, 3},
		},
		"nest": {
			yaml: strconv.Quote(strings.Trim(`
{{f <-}}:
  ids: |-
    {{f <-}}: 1
`, "\n")),
			ok: []interface{}{
				map[string]interface{}{
					"ids": []int{0, 1},
				},
			},
			ng: []interface{}{
				map[string]interface{}{
					"ids": []int{2, 3},
				},
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			var i interface{}
			if err := yaml.Unmarshal([]byte(tc.yaml), &i); err != nil {
				t.Fatalf("failed to unmarshal: %s", err)
			}
			v, err := template.Execute(i, map[string]interface{}{
				"f": leftArrowFunc(buildArg(assert.Contains)),
			})
			if err != nil {
				t.Fatalf("failed to execute: %s", err)
			}
			assertion := assert.Build(v)
			if err := assertion.Assert(tc.ok); err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if err := assertion.Assert(tc.ng); err == nil {
				t.Errorf("expected error but no error")
			}
		})
	}
}
