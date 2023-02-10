package assert

import (
	"github.com/bilus/scenarigo/errors"
)

// And returns a new assertion to ensure that value passes all assertions.
// If the assertions are empty, it returns an error.
func And(assertions ...Assertion) Assertion {
	return AssertionFunc(func(v interface{}) error {
		if len(assertions) == 0 {
			return errors.New("empty assertion list")
		}
		errs := []error{}
		for _, assertion := range assertions {
			assertion := assertion
			err := assertion.Assert(v)
			if err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) == 0 {
			return nil
		}
		return errors.Errors(errs...)
	})
}

// Or returns new assertion to ensure that value passes at least one of assertions.
// If the assertions are empty, it returns an error.
func Or(assertions ...Assertion) Assertion {
	return AssertionFunc(func(v interface{}) error {
		if len(assertions) == 0 {
			return errors.New("empty assertion list")
		}
		errs := []error{}
		for _, assertion := range assertions {
			assertion := assertion
			err := assertion.Assert(v)
			if err == nil {
				return nil
			}
			errs = append(errs, err)
		}
		return errors.Wrap(errors.Errors(errs...), "all assertions failed")
	})
}
