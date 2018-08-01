package envcfg

import (
	"reflect"
)

// Default read and parse values from struct tag `default:"some value"`
func Default(v interface{}) error {
	s := reflect.ValueOf(v)
	if err := checkType(s); err != nil {

		return err
	}

	return applyDefault(s)
}

func applyDefault(s reflect.Value) error {
	t := s.Elem().Type()
	m, err := getTagMap(t, "default")
	if err != nil {

		return err
	}

	return applyMap(s, m)
}
