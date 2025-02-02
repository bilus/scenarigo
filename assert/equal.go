package assert

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"

	"github.com/bilus/scenarigo/errors"
)

var (
	m        sync.RWMutex
	equalers []Equaler
)

// Equaler is the interface for custom equaler.
type Equaler interface {
	// Equal checks two values are equal or not.
	// If the ok is true, the err should be used as result.
	Equal(expected, got interface{}) (ok bool, err error)
}

// RegisterCustomEqualer appends eq as a custom equaler.
// Registered equaler will be used when all default equalers judge two values are not equal.
func RegisterCustomEqualer(eq Equaler) {
	m.Lock()
	defer m.Unlock()
	equalers = append(equalers, eq)
}

// EqualerFunc is an adaptor to allow the use of ordinary functions as Equaler.
func EqualerFunc(eq func(interface{}, interface{}) (bool, error)) Equaler {
	return equaler(eq)
}

type equaler func(interface{}, interface{}) (bool, error)

// Equal implements Equaler interface.
func (eq equaler) Equal(expected, got interface{}) (bool, error) {
	return eq(expected, got)
}

// Equal returns an assertion to ensure a value equals the expected value.
func Equal(expected interface{}) Assertion {
	return AssertionFunc(func(v interface{}) error {
		if n, ok := v.(json.Number); ok {
			switch expected.(type) {
			case int, int8, int16, int32, int64,
				uint, uint8, uint16, uint32, uint64:
				i, err := n.Int64()
				if err == nil {
					v = i
				}
			case float32, float64:
				f, err := n.Float64()
				if err == nil {
					v = f
				}
			}
		}

		// Attempt to coerce body to string if raw string specified.
		if isString(expected) {
			s, err := toString(v)
			if err == nil {
				v = s
			}
			log.Printf("%q", s)
		}

		if reflect.DeepEqual(v, expected) {
			return nil
		}

		if isNil(v) && isNil(expected) {
			return nil
		}

		m.RLock()
		defer m.RUnlock()
		for _, eq := range equalers {
			ok, err := eq.Equal(expected, v)
			if ok {
				return err
			}
		}

		if t := reflect.TypeOf(v); t != reflect.TypeOf(expected) {
			// try type conversion
			converted, err := convert(expected, t)
			if err == nil {
				if reflect.DeepEqual(v, converted) {
					return nil
				}
			}
			return errors.Errorf("expected %T (%+v) but got %T (%+v)", expected, expected, v, v)
		}
		return errors.Errorf("expected %+v but got %+v", expected, v)
	})
}

func toString(v interface{}) (interface{}, error) {
	if s, ok := v.([]uint8); ok {
		return string(s), nil
	}
	if s, ok := v.(string); ok {
		return s, nil
	}
	return nil, errors.Errorf("failed to convert %v (type %T) to string", v, v)
}

func isString(v interface{}) bool {
	switch v.(type) {
	case string:
		return true
	default:
		return false
	}
}

func isNil(i interface{}) bool {
	defer func() {
		// return false if IsNil panics
		_ = recover()
	}()
	if i == nil {
		return true
	}
	return reflect.ValueOf(i).IsNil()
}
