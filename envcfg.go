package envcfg

import (
	"os"
	"strings"

	"github.com/vporoshok/reflector"
)

type readConfig struct {
	prefix     string
	useDefault bool
	overrides  map[string]string
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
func WithDefault(overrides map[string]string) Option {

	return funcOption(func(cfg *readConfig) {
		cfg.useDefault = true
		cfg.overrides = overrides
	})
}

// Read environment into struct
func Read(v interface{}, opts ...Option) error {
	cfg := readConfig{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	r := reflector.New(v)
	if cfg.useDefault {
		if err := applyDefault(r, cfg.overrides); err != nil {

			return err
		}
	}

	m := r.ExtractTags("envcfg", reflector.WithoutMinus())
	for k, v := range m {
		if len(v) == 0 {
			v = strings.ToUpper(strings.Join(SplitWords(k), "_"))
		}
		var ok bool
		m[k], ok = os.LookupEnv(cfg.prefix + v)
		if !ok && len(cfg.prefix) > 0 {
			m[k], ok = os.LookupEnv(v)
		}
		if !ok {
			delete(m, k)

			continue
		}
	}

	return r.Apply(m)
}

// Default read and parse values from struct tag `default:"some value"`
func Default(v interface{}, overrides map[string]string) error {
	r := reflector.New(v)

	return applyDefault(r, overrides)
}

func applyDefault(r reflector.Reflector, overrides map[string]string) error {
	m := r.ExtractTags("default", reflector.WithoutEmpty())
	for k, v := range overrides {
		m[k] = v
	}

	return r.Apply(m)
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

// nolint:gocyclo
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
		_ = word.WriteByte(c)

		switch {
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
				_ = word.WriteByte(c)
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
		_ = res.WriteByte(s[i])
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
