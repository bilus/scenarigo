package assert

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/bilus/scenarigo/errors"
)

func TestBuild(t *testing.T) {
	str := `
deps:
- name: scenarigo
  version:
    major: 1
    minor: 2
    patch: 3
  tags:
    - go
    - test`
	var in interface{}
	if err := yaml.NewDecoder(strings.NewReader(str), yaml.UseOrderedMap()).Decode(&in); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	qs := []string{
		".deps[0].name",
		".deps[0].version.major",
		".deps[0].version.minor",
		".deps[0].version.patch",
		".deps[0].tags[0]",
		".deps[0].tags[1]",
	}
	assertion := Build(in)

	type info struct {
		Deps []map[string]interface{} `yaml:"deps"`
	}

	t.Run("no assertion", func(t *testing.T) {
		assertion := Build(nil)
		v := info{}
		if err := assertion.Assert(v); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("compare", func(t *testing.T) {
		if err := Build(Greater(1)).Assert(2); err != nil {
			t.Fatal(err)
		}
		if err := Build(GreaterOrEqual(1)).Assert(1); err != nil {
			t.Fatal(err)
		}
		if err := Build(Less(2)).Assert(1); err != nil {
			t.Fatal(err)
		}
		if err := Build(LessOrEqual(1)).Assert(1); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("ok", func(t *testing.T) {
		v := info{
			Deps: []map[string]interface{}{
				{
					"name": "scenarigo",
					"version": map[string]int{
						"major": 1,
						"minor": 2,
						"patch": 3,
					},
					"tags": []string{"go", "test"},
				},
			},
		}
		if err := assertion.Assert(v); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("ng", func(t *testing.T) {
		v := info{
			Deps: []map[string]interface{}{
				{
					"name": "Ruby on Rails",
					"version": map[string]int{
						"major": 2,
						"minor": 3,
						"patch": 4,
					},
					"tags": []string{"ruby", "http"},
				},
			},
		}
		err := assertion.Assert(v)
		if err == nil {
			t.Fatalf("expected error but no error")
		}
		var mperr *errors.MultiPathError
		if ok := errors.As(err, &mperr); !ok {
			t.Fatalf("expected errors.MultiPathError: %s", err)
		}
		if got, expect := len(mperr.Errs), len(qs); got != expect {
			t.Fatalf("expected %d but got %d", expect, got)
		}
		for i, e := range mperr.Errs {
			q := qs[i]
			if !strings.Contains(e.Error(), q) {
				t.Errorf(`"%s" does not contain "%s"`, e.Error(), q)
			}
		}
	})
	t.Run("assert nil", func(t *testing.T) {
		err := assertion.Assert(nil)
		if err == nil {
			t.Fatalf("expected error but no error")
		}
		var mperr *errors.MultiPathError
		if ok := errors.As(err, &mperr); !ok {
			t.Fatalf("expected errors.MultiPathError: %s", err)
		}
		if got, expect := len(mperr.Errs), len(qs); got != expect {
			t.Fatalf("expected %d but got %d", expect, got)
		}
		for i, e := range mperr.Errs {
			q := qs[i]
			if !strings.Contains(e.Error(), q) {
				t.Errorf(`"%s" does not contain "%s"`, e.Error(), q)
			}
		}
	})
}
