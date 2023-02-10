package marshaler

import (
	"reflect"

	"github.com/bilus/scenarigo/internal/reflectutil"
)

func init() {
	if err := Register(&textMarshaler{}); err != nil {
		panic(err)
	}
}

type textMarshaler struct{}

// MediaType implements RequestMarshaler interface.
func (m *textMarshaler) MediaType() string {
	return "text/plain"
}

// Marshal implements RequestMarshaler interface.
func (m *textMarshaler) Marshal(v interface{}) ([]byte, error) {
	s, err := reflectutil.ConvertString(reflect.ValueOf(v))
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}
