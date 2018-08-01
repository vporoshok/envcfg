package envcfg

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

func checkType(s reflect.Value) error {
	if s.Kind() == reflect.Ptr && s.Elem().Kind() == reflect.Struct {

		return nil
	}

	return InvalidObjectType
}

func getTagMap(t reflect.Type, tagName string) (map[string]string, error) {
	res := make(map[string]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		res[f.Name] = f.Tag.Get(tagName)
	}

	return res, nil
}

func applyMap(s reflect.Value, m map[string]string) error {
	s = s.Elem()
	for k, v := range m {
		if v == "" {

			continue
		}
		f := s.FieldByName(k)
		if err := processValue(f, v); err != nil {

			return err
		}
	}

	return nil
}

func processValue(s reflect.Value, v string) error {
	t := s.Type()

	switch s.Kind() {
	default:

		return InvalidFieldType.New(s.Kind().String())

	case reflect.String:
		s.SetString(v)

		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if s.Kind() == reflect.Int64 && t.PkgPath() == "time" && t.Name() == "Duration" {
			d, err := time.ParseDuration(v)
			if err != nil {

				return err
			}
			s.SetInt(int64(d))

			return nil
		}

		d, err := strconv.ParseInt(v, 0, t.Bits())
		if err != nil {

			return err
		}

		s.SetInt(d)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d, err := strconv.ParseUint(v, 0, t.Bits())
		if err != nil {

			return err
		}
		s.SetUint(d)

	case reflect.Bool:
		d, err := strconv.ParseBool(v)
		if err != nil {

			return err
		}
		s.SetBool(d)

	case reflect.Float32, reflect.Float64:
		d, err := strconv.ParseFloat(v, t.Bits())
		if err != nil {

			return err
		}
		s.SetFloat(d)

	case reflect.Slice:
		vals := strings.Split(v, ",")
		sl := reflect.MakeSlice(t, len(vals), len(vals))
		for i, val := range vals {
			err := processValue(sl.Index(i), val)
			if err != nil {

				return err
			}
		}
		s.Set(sl)
	}

	return nil
}
