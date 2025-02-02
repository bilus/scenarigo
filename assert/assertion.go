// Package assert provides value assertions.
package assert

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/zoncoen/query-go"
	yamlextractor "github.com/zoncoen/query-go/extractor/yaml"

	"github.com/bilus/scenarigo/errors"
)

// Assertion implements value assertion.
type Assertion interface {
	Assert(v interface{}) error
}

// AssertionFunc is an adaptor to allow the use of ordinary functions as assertions.
type AssertionFunc func(v interface{}) error

// Assert asserts the v.
func (f AssertionFunc) Assert(v interface{}) error {
	return f(v)
}

// Build creates an assertion from Go value.
func Build(expect interface{}) Assertion {
	var assertions []Assertion
	if expect != nil {
		assertions = build(query.New(
			query.ExtractByStructTag("yaml", "json"),
			query.CustomExtractFunc(yamlextractor.MapSliceExtractFunc(false)),
		), expect)
	}
	return AssertionFunc(func(v interface{}) error {
		errs := []error{}
		for _, assertion := range assertions {
			assertion := assertion
			if err := assertion.Assert(v); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			if len(errs) == 1 {
				return errs[0]
			}
			return errors.Errors(errs...)
		}
		return nil
	})
}

func build(q *query.Query, expect interface{}) []Assertion {
	var assertions []Assertion
	switch v := expect.(type) {
	case yaml.MapSlice:
		for _, item := range v {
			item := item
			key := fmt.Sprintf("%s", item.Key)
			assertions = append(assertions, build(q.Key(key), item.Value)...)
		}
	case []interface{}:
		for i, elm := range v {
			elm := elm
			assertions = append(assertions, build(q.Index(i), elm)...)
		}
	default:
		switch v := expect.(type) {
		case Assertion:
			assertions = append(assertions, AssertionFunc(func(val interface{}) error {
				got, err := q.Extract(val)
				if err != nil {
					return err
				}
				if err := v.Assert(got); err != nil {
					return errors.WithQuery(err, q)
				}
				return nil
			}))
		case func(*query.Query) Assertion:
			assertions = append(assertions, v(q))
		default:
			assertions = append(assertions, build(q, Equal(v))...)
		}
	}
	return assertions
}
