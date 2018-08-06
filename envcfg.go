package envcfg

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// InvalidObjectType returned if object passed to processing is not pointer to struct
	InvalidObjectType constantError = "expected pointer to struct"
	// InvalidFieldType returned if some field has unsupported type
	InvalidFieldType constantError = "unsupported type"
)

type constantError string

func (ce constantError) Error() string {

	return string(ce)
}

func (ce constantError) New(msg string) error {

	return causedError{
		err: ce,
		msg: msg,
	}
}

type causedError struct {
	err error
	msg string
}

func (ce causedError) Error() string {

	return fmt.Sprintf("%s: %s", ce.msg, ce.err.Error())
}

func (ce causedError) Cause() error {

	return ce.err
}

type readConfig struct {
	prefix     string
	useDefault bool
}

// Option to configure environment read
type Option interface {
	apply(*readConfig)
}

type funcOption func(*readConfig)

func (fo funcOption) apply(cfg *readConfig) {
	fo(cfg)
}

// WithPrefix add prefix to read variable
//
// If PREFIX_VARNAME is not setted, try to read VARNAME without prefix
func WithPrefix(prefix string) Option {

	return funcOption(func(cfg *readConfig) {
		cfg.prefix = prefix
	})
}

// WithDefault add read default values from tags before read environment
func WithDefault() Option {

	return funcOption(func(cfg *readConfig) {
		cfg.useDefault = true
	})
}

// Read environment into struct
func Read(v interface{}, opts ...Option) error {
	cfg := readConfig{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	s := reflect.ValueOf(v)
	if err := checkType(s); err != nil {

		return err
	}

	if cfg.useDefault {
		if err := applyDefault(s); err != nil {

			return err
		}
	}

	m := getTagMap(s.Elem().Type(), "envcfg")

	for k, v := range m {
		if v == "-" {
			delete(m, k)

			continue
		}
		if len(v) == 0 {
			v = strings.ToUpper(strings.Join(SplitWords(k), "_"))
		}
		m[k] = os.Getenv(cfg.prefix + v)
		if len(m[k]) == 0 && len(cfg.prefix) > 0 {
			m[k] = os.Getenv(v)
		}
		if len(m[k]) == 0 {
			delete(m, k)

			continue
		}
	}

	return applyMap(s, m)
}

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
	m := getTagMap(t, "default")

	return applyMap(s, m)
}

func checkType(s reflect.Value) error {
	if s.Kind() == reflect.Ptr && s.Elem().Kind() == reflect.Struct {

		return nil
	}

	return InvalidObjectType
}

func getTagMap(t reflect.Type, tagName string) map[string]string {
	res := make(map[string]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		res[f.Name] = f.Tag.Get(tagName)
	}

	return res
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

// SplitWords split string to words
//
// JSON42_File42An42dSO42ME_More -> JSON42, File42, An42d, SO42ME, More
func SplitWords(s string) []string {
	var res []string

	chunks := strings.Split(s, "_")
	for _, chunk := range chunks {
		res = append(res, splitChunk(chunk)...)
	}

	return res
}

func splitChunk(chunk string) []string {
	const (
		unknown = iota
		lower
		number
		upper
	)

	isUpper := func(c byte) bool { return c >= 'A' && c <= 'Z' }
	isNumber := func(c byte) bool { return c >= '0' && c <= '9' }

	if len(chunk) == 0 {

		return nil
	}
	var res []string
	word := &strings.Builder{}
	prev := unknown
	for i := len(chunk) - 1; i >= 0; i-- {
		c := chunk[i]
		word.WriteByte(c)

		switch true {
		case isUpper(c):
			if prev == lower {
				res = append(res, reverseString(word.String()))
				word.Reset()
				prev = unknown

				continue
			}

			prev = upper

		case isNumber(c):
			prev = number

		default:
			if prev == upper {
				w := word.String()
				res = append(res, reverseString(w[:len(w)-1]))
				word.Reset()
				word.WriteByte(c)
			}

			prev = lower
		}
	}

	if word.Len() > 0 {
		res = append(res, reverseString(word.String()))
	}

	return reverseSlice(res)
}

func reverseString(s string) string {
	res := &strings.Builder{}
	for i := len(s) - 1; i >= 0; i-- {
		res.WriteByte(s[i])
	}

	return res.String()
}

func reverseSlice(s []string) []string {
	res := make([]string, len(s))
	for i := range res {
		res[i] = s[len(s)-i-1]
	}

	return res
}
