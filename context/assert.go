package context

import (
	"github.com/pkg/errors"

	"github.com/bilus/scenarigo/assert"
)

var assertions = map[string]interface{}{
	"and":                listArgsLeftArrowFunc(buildArgs(assert.And)),
	"or":                 listArgsLeftArrowFunc(buildArgs(assert.Or)),
	"notZero":            assert.NotZero(),
	"contains":           leftArrowFunc(buildArg(assert.Contains)),
	"notContains":        leftArrowFunc(buildArg(assert.NotContains)),
	"regexp":             assert.Regexp,
	"greaterThan":        assert.Greater,
	"greaterThanOrEqual": assert.GreaterOrEqual,
	"lessThan":           assert.Less,
	"lessThanOrEqual":    assert.LessOrEqual,
	"length":             assert.Length,
}

func buildArg(base func(assert.Assertion) assert.Assertion) func(interface{}) assert.Assertion {
	return func(arg interface{}) assert.Assertion {
		assertion, ok := arg.(assert.Assertion)
		if !ok {
			assertion = assert.Build(arg)
		}
		return base(assertion)
	}
}

type leftArrowFunc func(interface{}) assert.Assertion

func (f leftArrowFunc) Exec(arg interface{}) (interface{}, error) {
	assertion, ok := arg.(assert.Assertion)
	if !ok {
		return nil, errors.New("argument must be a assert.Assertion")
	}
	return f(assertion), nil
}

func (leftArrowFunc) UnmarshalArg(unmarshal func(interface{}) error) (interface{}, error) {
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return nil, err
	}
	return assert.Build(i), nil
}

func buildArgs(base func(...assert.Assertion) assert.Assertion) func(...interface{}) assert.Assertion {
	return func(args ...interface{}) assert.Assertion {
		var assertions []assert.Assertion
		for _, arg := range args {
			arg := arg
			assertion, ok := arg.(assert.Assertion)
			if !ok {
				assertion = assert.Build(arg)
			}
			assertions = append(assertions, assertion)
		}
		return base(assertions...)
	}
}

type listArgsLeftArrowFunc func(args ...interface{}) assert.Assertion

func (f listArgsLeftArrowFunc) Exec(arg interface{}) (interface{}, error) {
	assertions, ok := arg.([]interface{})
	if !ok {
		return nil, errors.New("argument must be a slice of interface{}")
	}
	return f(assertions...), nil
}

func (listArgsLeftArrowFunc) UnmarshalArg(unmarshal func(interface{}) error) (interface{}, error) {
	var args []interface{}
	if err := unmarshal(&args); err != nil {
		return nil, err
	}
	return args, nil
}
