package assert

import (
	"reflect"

	"github.com/bilus/scenarigo/errors"
	"github.com/bilus/scenarigo/internal/reflectutil"
)

// Contains returns an assertion to ensure a value contains the value.
func Contains(assertion Assertion) Assertion {
	return AssertionFunc(func(v interface{}) error {
		vv, err := arrayOrSlice(v)
		if err != nil {
			return err
		}
		if err := contains(assertion, vv); err != nil {
			return errors.Wrap(err, "doesn't contain expected value")
		}
		return nil
	})
}

// NotContains returns an assertion to ensure a value doesn't contain the value.
func NotContains(assertion Assertion) Assertion {
	return AssertionFunc(func(v interface{}) error {
		vv, err := arrayOrSlice(v)
		if err != nil {
			return err
		}
		if err := contains(assertion, vv); err == nil {
			return errors.Wrap(err, "contains the value")
		}
		return nil
	})
}

func arrayOrSlice(v interface{}) (reflect.Value, error) {
	vv := reflectutil.Elem(reflect.ValueOf(v))
	switch vv.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		return reflect.Value{}, errors.New("expected an array")
	}
	return vv, nil
}

func contains(assertion Assertion, v reflect.Value) error {
	if v.Len() == 0 {
		return errors.New("empty")
	}
	var err error
	for i := 0; i < v.Len(); i++ {
		e := v.Index(i).Interface()
		if err = assertion.Assert(e); err == nil {
			return nil
		}
	}
	return errors.Wrap(err, "last error")
}
